package notification

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/service"
)

// PublishNotificationUseCase defines the interface for publish notifications.
type PublishNotificationUseCase interface {
	// Execute enqueues a notification job for a buyer.
	Execute(ctx context.Context, buyerID int, payload any) error
}

type publishNotificationUseCase struct {
	pushQueue service.PushNotificationQueue
}

var _ PublishNotificationUseCase = (*publishNotificationUseCase)(nil)

// NewPublishNotificationUseCase creates a new instance of PublishNotificationUseCase.
func NewPublishNotificationUseCase(
	pushQueue service.PushNotificationQueue,
) PublishNotificationUseCase {
	return &publishNotificationUseCase{
		pushQueue: pushQueue,
	}
}

func (uc *publishNotificationUseCase) Execute(ctx context.Context, buyerID int, payload any) error {
	return uc.pushQueue.Enqueue(ctx, buyerID, payload)
}
