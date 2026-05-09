package relay

import (
	"context"
	"log"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

// OutboxRelay polls the outbox table and forwards messages to SQS.
// Uses a 3-phase approach to avoid holding DB locks during external I/O:
//
//	Phase 1 (TX1): Claim pending messages (status → processing)
//	Phase 2 (no TX): Send to SQS
//	Phase 3 (TX2): Mark as processed or record failure with backoff
//
// 複数インスタンス起動時の安全性:
//   - Claim は FOR UPDATE SKIP LOCKED により同一行を二重 claim しない
//   - MarkProcessed / MarkFailed は claimed_by フィルタで自インスタンス分のみ更新
//   - クラッシュ等で取り残された処理中行は RecoverStale により pending 復元
//
// SQS への enqueue は at-least-once 配送が前提。worker 側ハンドラは冪等であること。
type OutboxRelay struct {
	outboxRepo repository.OutboxRepository
	jobQueue   service.JobQueue
	interval   time.Duration
	batchSize  int
	instanceID string
}

// NewOutboxRelay creates a new OutboxRelay.
func NewOutboxRelay(
	outboxRepo repository.OutboxRepository,
	jobQueue service.JobQueue,
	interval time.Duration,
	batchSize int,
	instanceID string,
) *OutboxRelay {
	return &OutboxRelay{
		outboxRepo: outboxRepo,
		jobQueue:   jobQueue,
		interval:   interval,
		batchSize:  batchSize,
		instanceID: instanceID,
	}
}

// Run drives the polling loop until ctx is canceled.
func (r *OutboxRelay) Run(ctx context.Context) {
	log.Printf("OutboxRelay: started (interval=%s, batch=%d, instance=%s)", r.interval, r.batchSize, r.instanceID)
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("OutboxRelay: stopping")
			return
		case <-ticker.C:
			r.relay(ctx)
		}
	}
}

func (r *OutboxRelay) relay(ctx context.Context) {
	// Phase 1: Claim pending messages (DB only, commits immediately)
	msgs, err := r.outboxRepo.Claim(ctx, r.batchSize, r.instanceID)
	if err != nil {
		log.Printf("OutboxRelay: claim error: %v", err)
		return
	}
	if len(msgs) == 0 {
		return
	}

	log.Printf("OutboxRelay: claimed %d messages", len(msgs))

	// Phase 2: Send to SQS (outside any transaction)
	var successIDs []int64
	for _, msg := range msgs {
		if err := r.jobQueue.Enqueue(ctx, msg.JobType, msg.Payload); err != nil {
			// Phase 3a: Record failure with backoff
			if markErr := r.outboxRepo.MarkFailed(ctx, msg.ID, err.Error(), r.instanceID); markErr != nil {
				log.Printf("OutboxRelay: failed to mark message %d as failed: %v", msg.ID, markErr)
			}
			log.Printf("OutboxRelay: failed to enqueue message %d: %v", msg.ID, err)
			continue
		}
		successIDs = append(successIDs, msg.ID)
		log.Printf("OutboxRelay: successfully enqueued message %d (type: %s)", msg.ID, msg.JobType)
	}

	// Phase 3b: Mark successful messages as processed
	if len(successIDs) > 0 {
		if err := r.outboxRepo.MarkProcessed(ctx, successIDs, r.instanceID); err != nil {
			log.Printf("OutboxRelay: failed to mark %d messages as processed: %v", len(successIDs), err)
		}
	}
}
