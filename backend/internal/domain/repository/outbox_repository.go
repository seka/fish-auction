package repository

import (
	"context"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// OutboxRepository manages outbox messages for the transactional outbox pattern.
type OutboxRepository interface {
	// InsertEmailJob serializes and inserts an email job.
	InsertEmailJob(ctx context.Context, to string, resetURL string, emailType string) error

	// InsertPushNotificationJob serializes and inserts a push notification job.
	InsertPushNotificationJob(ctx context.Context, buyerID int, payload any) error

	// Claim claims pending messages for processing.
	Claim(ctx context.Context, batchSize int, instanceID string) ([]*model.OutboxMessage, error)

	// MarkProcessed sets status to processed for successfully sent messages.
	MarkProcessed(ctx context.Context, ids []int64) error

	// MarkFailed records a send failure with exponential backoff.
	// When max_attempts is reached, status becomes failed (poison message isolation).
	MarkFailed(ctx context.Context, id int64, lastError string) error

	// RecoverStale resets messages stuck in processing state back to pending.
	RecoverStale(ctx context.Context, timeout time.Duration) (int64, error)

	// DeleteProcessedBefore removes processed messages older than the given time.
	DeleteProcessedBefore(ctx context.Context, before time.Time) (int64, error)
}
