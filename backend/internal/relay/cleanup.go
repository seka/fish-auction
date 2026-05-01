package relay

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// OutboxCleanup periodically removes old processed outbox records
// and recovers messages stuck in processing state.
type OutboxCleanup struct {
	outboxRepo    repository.OutboxRepository
	retention     time.Duration
	cleanInterval time.Duration
	staleTimeout  time.Duration
	staleInterval time.Duration
}

// NewOutboxCleanup creates a new OutboxCleanup.
func NewOutboxCleanup(
	outboxRepo repository.OutboxRepository,
	retention time.Duration,
	cleanInterval time.Duration,
	staleTimeout time.Duration,
	staleInterval time.Duration,
) *OutboxCleanup {
	return &OutboxCleanup{
		outboxRepo:    outboxRepo,
		retention:     retention,
		cleanInterval: cleanInterval,
		staleTimeout:  staleTimeout,
		staleInterval: staleInterval,
	}
}

// Start begins the cleanup and stale recovery loops as goroutines.
func (c *OutboxCleanup) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(2)

	// Processed record cleanup
	go func() {
		defer wg.Done()
		log.Printf("OutboxCleanup: started (retention=%s, interval=%s)", c.retention, c.cleanInterval)
		ticker := time.NewTicker(c.cleanInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("OutboxCleanup: stopping")
				return
			case <-ticker.C:
				before := time.Now().Add(-c.retention)
				deleted, err := c.outboxRepo.DeleteProcessedBefore(ctx, before)
				if err != nil {
					log.Printf("OutboxCleanup: error: %v", err)
				} else if deleted > 0 {
					log.Printf("OutboxCleanup: deleted %d processed messages", deleted)
				}
			}
		}
	}()

	// Stale message recovery
	go func() {
		defer wg.Done()
		log.Printf("OutboxCleanup: stale recovery started (timeout=%s, interval=%s)", c.staleTimeout, c.staleInterval)
		ticker := time.NewTicker(c.staleInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("OutboxCleanup: stale recovery stopping")
				return
			case <-ticker.C:
				recovered, err := c.outboxRepo.RecoverStale(ctx, c.staleTimeout)
				if err != nil {
					log.Printf("OutboxCleanup: stale recovery error: %v", err)
				} else if recovered > 0 {
					log.Printf("OutboxCleanup: recovered %d stale messages", recovered)
				}
			}
		}
	}()
}
