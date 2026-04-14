package service

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// JobQueue defines the interface for asynchronous job messaging.
type JobQueue interface {
	Enqueue(ctx context.Context, jobType model.JobType, payload []byte) error
}
