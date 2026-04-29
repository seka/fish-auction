package sqs

import (
	"context"
	"errors"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue"
	notificationMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

type sqsPushNotificationService struct {
	jobQueue queue.JobQueue
}

var _ service.PushNotificationService = (*sqsPushNotificationService)(nil)

// NewPushNotificationService creates a PushNotificationService that enqueues jobs via SQS.
func NewPushNotificationService(jobQueue queue.JobQueue) service.PushNotificationService {
	return &sqsPushNotificationService{jobQueue: jobQueue}
}

func (s *sqsPushNotificationService) Send(ctx context.Context, sub *model.PushSubscription, payload any) error {
	return errors.New("Send is not supported in sqsPushNotificationService")
}

func (s *sqsPushNotificationService) PublishToBuyer(ctx context.Context, buyerID int, payload any) error {
	jobParams := notificationMessage.PushNotificationMessage{
		BuyerID: buyerID,
		Payload: payload,
	}
	return s.jobQueue.Enqueue(ctx, model.JobTypePushNotification, jobParams)
}
