package testing

import (
	"context"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// MockOutboxRepository is a mock implementation of OutboxRepository for testing.
type MockOutboxRepository struct {
	InsertEmailJobFunc        func(ctx context.Context, to string, resetURL string, emailType string) error
	InsertPushJobFunc         func(ctx context.Context, jobType model.JobType, buyerID int, title, body, url string) error
	ClaimFunc                 func(ctx context.Context, batchSize int, instanceID string) ([]*model.OutboxMessage, error)
	MarkProcessedFunc         func(ctx context.Context, ids []int64) error
	MarkFailedFunc            func(ctx context.Context, id int64, lastError string) error
	RecoverStaleFunc          func(ctx context.Context, timeout time.Duration) (int64, error)
	DeleteProcessedBeforeFunc func(ctx context.Context, before time.Time) (int64, error)
}

var _ repository.OutboxRepository = (*MockOutboxRepository)(nil)

func (m *MockOutboxRepository) InsertEmailJob(ctx context.Context, to, resetURL, emailType string) error {
	if m.InsertEmailJobFunc != nil {
		return m.InsertEmailJobFunc(ctx, to, resetURL, emailType)
	}
	return nil
}

func (m *MockOutboxRepository) InsertPushJob(ctx context.Context, jobType model.JobType, buyerID int, title, body, url string) error {
	if m.InsertPushJobFunc != nil {
		return m.InsertPushJobFunc(ctx, jobType, buyerID, title, body, url)
	}
	return nil
}

func (m *MockOutboxRepository) Claim(ctx context.Context, batchSize int, instanceID string) ([]*model.OutboxMessage, error) {
	if m.ClaimFunc != nil {
		return m.ClaimFunc(ctx, batchSize, instanceID)
	}
	return nil, nil
}

func (m *MockOutboxRepository) MarkProcessed(ctx context.Context, ids []int64) error {
	if m.MarkProcessedFunc != nil {
		return m.MarkProcessedFunc(ctx, ids)
	}
	return nil
}

func (m *MockOutboxRepository) MarkFailed(ctx context.Context, id int64, lastError string) error {
	if m.MarkFailedFunc != nil {
		return m.MarkFailedFunc(ctx, id, lastError)
	}
	return nil
}

func (m *MockOutboxRepository) RecoverStale(ctx context.Context, timeout time.Duration) (int64, error) {
	if m.RecoverStaleFunc != nil {
		return m.RecoverStaleFunc(ctx, timeout)
	}
	return 0, nil
}

func (m *MockOutboxRepository) DeleteProcessedBefore(ctx context.Context, before time.Time) (int64, error) {
	if m.DeleteProcessedBeforeFunc != nil {
		return m.DeleteProcessedBeforeFunc(ctx, before)
	}
	return 0, nil
}
