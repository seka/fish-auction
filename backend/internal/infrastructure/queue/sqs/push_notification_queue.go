package sqs

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue"
	notificationMessage "github.com/seka/fish-auction/backend/internal/worker/message"
)

type pushNotificationQueue struct {
	jobQueue queue.JobQueue
}

var _ service.PushNotificationQueue = (*pushNotificationQueue)(nil)

// NewPushNotificationQueue creates a PushNotificationQueue that enqueues jobs via SQS.
func NewPushNotificationQueue(jobQueue queue.JobQueue) service.PushNotificationQueue {
	return &pushNotificationQueue{jobQueue: jobQueue}
}

func (s *pushNotificationQueue) Enqueue(ctx context.Context, buyerID int, payload any) error {
	jobParams := notificationMessage.PushNotificationMessage{
		BuyerID: buyerID,
		Payload: payload,
	}
	return s.jobQueue.Enqueue(ctx, model.JobTypePushNotification, jobParams)
}

func (s *pushNotificationQueue) Dequeue(ctx context.Context, waitTimeSeconds int32) ([]*model.JobMessage, error) {
	return s.jobQueue.Dequeue(ctx, waitTimeSeconds)
}

func (s *pushNotificationQueue) DeleteMessage(ctx context.Context, message *model.JobMessage) error {
	return s.jobQueue.DeleteMessage(ctx, message)
}
