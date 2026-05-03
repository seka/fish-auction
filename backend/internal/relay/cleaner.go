package relay

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// OutboxCleaner periodically removes old processed outbox records
// and recovers messages stuck in processing state.
type OutboxCleaner struct {
	outboxRepo    repository.OutboxRepository
	retention     time.Duration
	cleanInterval time.Duration
	staleTimeout  time.Duration
	staleInterval time.Duration
}

// NewOutboxCleaner creates a new OutboxCleaner.
func NewOutboxCleaner(
	outboxRepo repository.OutboxRepository,
	retention time.Duration,
	cleanInterval time.Duration,
	staleTimeout time.Duration,
	staleInterval time.Duration,
) *OutboxCleaner {
	return &OutboxCleaner{
		outboxRepo:    outboxRepo,
		retention:     retention,
		cleanInterval: cleanInterval,
		staleTimeout:  staleTimeout,
		staleInterval: staleInterval,
	}
}

// Run drives the cleanup and stale-recovery loops until ctx is canceled.
func (c *OutboxCleaner) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		c.cleanLoop(ctx)
	}()

	go func() {
		defer wg.Done()
		c.staleLoop(ctx)
	}()

	wg.Wait()
}

func (c *OutboxCleaner) cleanLoop(ctx context.Context) {
	log.Printf("OutboxCleaner: started (retention=%s, interval=%s)", c.retention, c.cleanInterval)
	ticker := time.NewTicker(c.cleanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("OutboxCleaner: stopping")
			return
		case <-ticker.C:
			before := time.Now().Add(-c.retention)
			deleted, err := c.outboxRepo.DeleteProcessedBefore(ctx, before)
			if err != nil {
				log.Printf("OutboxCleaner: error: %v", err)
			} else if deleted > 0 {
				log.Printf("OutboxCleaner: deleted %d processed messages", deleted)
			}
		}
	}
}

func (c *OutboxCleaner) staleLoop(ctx context.Context) {
	log.Printf("OutboxCleaner: stale recovery started (timeout=%s, interval=%s)", c.staleTimeout, c.staleInterval)
	ticker := time.NewTicker(c.staleInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("OutboxCleaner: stale recovery stopping")
			return
		case <-ticker.C:
			recovered, err := c.outboxRepo.RecoverStale(ctx, c.staleTimeout)
			if err != nil {
				log.Printf("OutboxCleaner: stale recovery error: %v", err)
			} else if recovered > 0 {
				log.Printf("OutboxCleaner: recovered %d stale messages", recovered)
			}
		}
	}
}
