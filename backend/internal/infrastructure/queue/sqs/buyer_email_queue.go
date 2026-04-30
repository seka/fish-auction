package sqs

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue"

	emailMessage "github.com/seka/fish-auction/backend/internal/worker/message"
)

type buyerEmailQueue struct {
	jobQueue queue.JobQueue
}

var _ service.BuyerEmailQueue = (*buyerEmailQueue)(nil)

// NewBuyerEmailQueue creates a BuyerEmailQueue that enqueues email jobs via SQS.
func NewBuyerEmailQueue(jobQueue queue.JobQueue) service.BuyerEmailQueue {
	return &buyerEmailQueue{jobQueue: jobQueue}
}

func (s *buyerEmailQueue) EnqueueBuyerPasswordReset(ctx context.Context, to, resetURL string) error {
	wire := emailMessage.EmailMessage{
		EmailType: emailMessage.EmailTypeBuyerPasswordReset,
		To:        to,
		ResetURL:  resetURL,
	}
	return s.jobQueue.Enqueue(ctx, model.JobTypeEmail, wire)
}

func (s *buyerEmailQueue) Dequeue(ctx context.Context, waitTimeSeconds int32) ([]*model.JobMessage, error) {
	return s.jobQueue.Dequeue(ctx, waitTimeSeconds)
}

func (s *buyerEmailQueue) DeleteMessage(ctx context.Context, message *model.JobMessage) error {
	return s.jobQueue.DeleteMessage(ctx, message)
}
