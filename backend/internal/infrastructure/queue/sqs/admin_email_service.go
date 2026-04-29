package sqs

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue"

	emailMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

type adminEmailService struct {
	jobQueue queue.JobQueue
}

var _ service.AdminEmailService = (*adminEmailService)(nil)

// NewAdminEmailService creates an AdminEmailService that enqueues email jobs via SQS.
func NewAdminEmailService(jobQueue queue.JobQueue) service.AdminEmailService {
	return &adminEmailService{jobQueue: jobQueue}
}

func (s *adminEmailService) SendAdminPasswordReset(ctx context.Context, to, resetURL string) error {
	wire := emailMessage.EmailMessage{
		EmailType: emailMessage.EmailTypeAdminPasswordReset,
		To:        to,
		ResetURL:  resetURL,
	}
	return s.jobQueue.Enqueue(ctx, model.JobTypeEmail, wire)
}
