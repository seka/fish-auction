package mock

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue"
)

// MockJobQueue is a test double for queue.JobQueue.
type MockJobQueue struct {
	EnqueueFunc func(ctx context.Context, jobType model.JobType, payload any) error
}

var _ queue.JobQueue = (*MockJobQueue)(nil)

func (m *MockJobQueue) Enqueue(ctx context.Context, jobType model.JobType, payload any) error {
	return m.EnqueueFunc(ctx, jobType, payload)
}

func (m *MockJobQueue) Dequeue(_ context.Context, _ int32) ([]*model.JobMessage, error) {
	return nil, nil
}

func (m *MockJobQueue) DeleteMessage(_ context.Context, _ *model.JobMessage) error {
	return nil
}
