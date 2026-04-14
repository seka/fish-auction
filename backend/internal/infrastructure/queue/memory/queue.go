package memory

import (
	"context"
	"sync"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

// Job represents an in-memory job.
type Job struct {
	JobType model.JobType
	Payload []byte
}

// MemoryQueue implements service.JobQueue in-memory for testing.
type MemoryQueue struct {
	mu   sync.Mutex
	Jobs []Job
}

var _ service.JobQueue = (*MemoryQueue)(nil)

// NewMemoryQueue creates a new in-memory job queue.
func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		Jobs: make([]Job, 0),
	}
}

// Enqueue adds a job to the in-memory queue.
func (q *MemoryQueue) Enqueue(ctx context.Context, jobType model.JobType, payload []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Jobs = append(q.Jobs, Job{
		JobType: jobType,
		Payload: payload,
	})
	return nil
}
