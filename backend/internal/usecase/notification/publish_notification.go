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
	pushSvc service.PushNotificationService
}

var _ PublishNotificationUseCase = (*publishNotificationUseCase)(nil)

// NewPublishNotificationUseCase creates a new instance of PublishNotificationUseCase.
func NewPublishNotificationUseCase(
	pushSvc service.PushNotificationService,
) PublishNotificationUseCase {
	return &publishNotificationUseCase{
		pushSvc: pushSvc,
	}
}

func (uc *publishNotificationUseCase) Execute(ctx context.Context, buyerID int, payload any) error {
	return uc.pushSvc.PublishToBuyer(ctx, buyerID, payload)
}
