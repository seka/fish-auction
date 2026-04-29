package sqs

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue"

	emailMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

type sqsAdminEmailService struct {
	jobQueue queue.JobQueue
}

var _ service.AdminEmailService = (*sqsAdminEmailService)(nil)

// NewAdminEmailService creates an AdminEmailService that enqueues email jobs via SQS.
func NewAdminEmailService(jobQueue queue.JobQueue) service.AdminEmailService {
	return &sqsAdminEmailService{jobQueue: jobQueue}
}

func (s *sqsAdminEmailService) SendAdminPasswordReset(ctx context.Context, to, resetURL string) error {
	wire := emailMessage.EmailMessage{
		EmailType: emailMessage.EmailTypeAdminPasswordReset,
		To:        to,
		ResetURL:  resetURL,
	}
	return s.jobQueue.Enqueue(ctx, model.JobTypeEmail, wire)
}
