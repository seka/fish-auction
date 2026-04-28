package notification

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	emailMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

// PublishEmailUseCase defines the interface for enqueuing email jobs.
type PublishEmailUseCase interface {
	Execute(ctx context.Context, msg emailMessage.EmailMessage) error
}

type publishEmailUseCase struct {
	jobQueue service.JobQueue
}

var _ PublishEmailUseCase = (*publishEmailUseCase)(nil)

// NewPublishEmailUseCase creates a new instance of PublishEmailUseCase.
func NewPublishEmailUseCase(jobQueue service.JobQueue) PublishEmailUseCase {
	return &publishEmailUseCase{jobQueue: jobQueue}
}

func (uc *publishEmailUseCase) Execute(ctx context.Context, msg emailMessage.EmailMessage) error {
	return uc.jobQueue.Enqueue(ctx, model.JobTypeEmail, msg)
}
