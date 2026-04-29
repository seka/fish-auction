package sqs

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue"

	emailMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

type buyerEmailService struct {
	jobQueue queue.JobQueue
}

var _ service.BuyerEmailService = (*buyerEmailService)(nil)

// NewBuyerEmailService creates a BuyerEmailService that enqueues email jobs via SQS.
func NewBuyerEmailService(jobQueue queue.JobQueue) service.BuyerEmailService {
	return &buyerEmailService{jobQueue: jobQueue}
}

func (s *buyerEmailService) SendBuyerPasswordReset(ctx context.Context, to, resetURL string) error {
	wire := emailMessage.EmailMessage{
		EmailType: emailMessage.EmailTypeBuyerPasswordReset,
		To:        to,
		ResetURL:  resetURL,
	}
	return s.jobQueue.Enqueue(ctx, model.JobTypeEmail, wire)
}
