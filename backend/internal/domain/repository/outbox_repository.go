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

	// InsertPushJob serializes and inserts a push notification job.
	// jobType must be one of JobTypePush* values; title/body/url are delivered as-is to the browser Service Worker.
	InsertPushJob(ctx context.Context, jobType model.JobType, buyerID int, title, body, url string) error

	// Claim claims pending messages for processing.
	Claim(ctx context.Context, batchSize int, instanceID string) ([]*model.OutboxMessage, error)

	// MarkProcessed sets status to processed for successfully sent messages.
	// claimedBy 所有者チェック: 別インスタンスが奪取したメッセージは更新しない。
	MarkProcessed(ctx context.Context, ids []int64, claimedBy string) error

	// MarkFailed records a send failure with exponential backoff.
	// When max_attempts is reached, status becomes failed (poison message isolation).
	// claimedBy 所有者チェック: 別インスタンスが奪取したメッセージは更新しない。
	MarkFailed(ctx context.Context, id int64, lastError, claimedBy string) error

	// RecoverStale resets messages stuck in processing state back to pending.
	RecoverStale(ctx context.Context, timeout time.Duration) (int64, error)

	// DeleteProcessedBefore removes processed messages older than the given time.
	DeleteProcessedBefore(ctx context.Context, before time.Time) (int64, error)
}
