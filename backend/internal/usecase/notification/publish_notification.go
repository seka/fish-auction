package notification

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	notificationMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

// PublishNotificationUseCase defines the interface for publish notifications.
type PublishNotificationUseCase interface {
	// Execute enqueues a notification job for a buyer.
	Execute(ctx context.Context, buyerID int, payload any) error
}

type publishNotificationUseCase struct {
	jobQueue service.JobQueue
}

var _ PublishNotificationUseCase = (*publishNotificationUseCase)(nil)

// NewPublishNotificationUseCase creates a new instance of PublishNotificationUseCase.
func NewPublishNotificationUseCase(
	jobQueue service.JobQueue,
) PublishNotificationUseCase {
	return &publishNotificationUseCase{
		jobQueue: jobQueue,
	}
}

func (uc *publishNotificationUseCase) Execute(ctx context.Context, buyerID int, payload any) error {
	jobParams := notificationMessage.PushNotificationMessage{
		BuyerID: buyerID,
		Payload: payload,
	}

	return uc.jobQueue.Enqueue(ctx, model.JobTypePushNotification, jobParams)
}
