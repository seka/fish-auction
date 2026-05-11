package relay

import (
	"context"
	"log/slog"
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
type OutboxRelay struct {
	outboxRepo repository.OutboxRepository
	jobQueue   service.JobQueue
	interval   time.Duration
	batchSize  int
	instanceID string
	logger     *slog.Logger
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
		logger:     slog.With("component", "outbox_relay", "instance_id", instanceID),
	}
}

// Run drives the polling loop until ctx is canceled.
func (r *OutboxRelay) Run(ctx context.Context) {
	r.logger.Info("outbox relay started", "interval", r.interval.String(), "batch", r.batchSize)
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			r.logger.Info("outbox relay stopping")
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
		r.logger.Error("claim error", "err", err)
		return
	}
	if len(msgs) == 0 {
		return
	}

	r.logger.Info("claimed messages", "count", len(msgs))

	// Phase 2: Send to SQS (outside any transaction)
	var successIDs []int64
	for _, msg := range msgs {
		if err := r.jobQueue.Enqueue(ctx, msg.JobType, msg.Payload); err != nil {
			// Phase 3a: Record failure with backoff
			if markErr := r.outboxRepo.MarkFailed(ctx, msg.ID, err.Error(), r.instanceID); markErr != nil {
				r.logger.Error("failed to mark message as failed", "message_id", msg.ID, "err", markErr)
			}
			r.logger.Error("failed to enqueue message", "message_id", msg.ID, "err", err)
			continue
		}
		successIDs = append(successIDs, msg.ID)
		r.logger.Info("successfully enqueued message", "message_id", msg.ID, "job_type", msg.JobType)
	}

	// Phase 3b: Mark successful messages as processed
	if len(successIDs) > 0 {
		if err := r.outboxRepo.MarkProcessed(ctx, successIDs, r.instanceID); err != nil {
			r.logger.Error("failed to mark messages as processed", "count", len(successIDs), "err", err)
		}
	}
}
