package sqs

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue"

	emailMessage "github.com/seka/fish-auction/backend/internal/worker/message"
)

type adminEmailQueue struct {
	jobQueue queue.JobQueue
}

var _ service.AdminEmailQueue = (*adminEmailQueue)(nil)

// NewAdminEmailQueue creates an AdminEmailQueue that enqueues email jobs via SQS.
func NewAdminEmailQueue(jobQueue queue.JobQueue) service.AdminEmailQueue {
	return &adminEmailQueue{jobQueue: jobQueue}
}

func (s *adminEmailQueue) EnqueueAdminPasswordReset(ctx context.Context, to, resetURL string) error {
	wire := emailMessage.EmailMessage{
		EmailType: emailMessage.EmailTypeAdminPasswordReset,
		To:        to,
		ResetURL:  resetURL,
	}
	return s.jobQueue.Enqueue(ctx, model.JobTypeEmail, wire)
}

func (s *adminEmailQueue) Dequeue(ctx context.Context, waitTimeSeconds int32) ([]*model.JobMessage, error) {
	return s.jobQueue.Dequeue(ctx, waitTimeSeconds)
}

func (s *adminEmailQueue) DeleteMessage(ctx context.Context, message *model.JobMessage) error {
	return s.jobQueue.DeleteMessage(ctx, message)
}
