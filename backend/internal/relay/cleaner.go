package relay

import (
	"context"
	"log/slog"
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
	logger        *slog.Logger
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
		logger:        slog.With("component", "outbox_cleaner"),
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
	c.logger.Info("clean loop started", "retention", c.retention.String(), "interval", c.cleanInterval.String())
	ticker := time.NewTicker(c.cleanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("clean loop stopping")
			return
		case <-ticker.C:
			before := time.Now().Add(-c.retention)
			deleted, err := c.outboxRepo.DeleteProcessedBefore(ctx, before)
			if err != nil {
				c.logger.Error("clean loop error", "err", err)
			} else if deleted > 0 {
				c.logger.Info("deleted processed messages", "count", deleted)
			}
		}
	}
}

func (c *OutboxCleaner) staleLoop(ctx context.Context) {
	c.logger.Info("stale recovery started", "timeout", c.staleTimeout.String(), "interval", c.staleInterval.String())
	ticker := time.NewTicker(c.staleInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stale recovery stopping")
			return
		case <-ticker.C:
			recovered, err := c.outboxRepo.RecoverStale(ctx, c.staleTimeout)
			if err != nil {
				c.logger.Error("stale recovery error", "err", err)
			} else if recovered > 0 {
				c.logger.Info("recovered stale messages", "count", recovered)
			}
		}
	}
}
