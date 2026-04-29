package push_notification

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type noopPushNotificationService struct{}

func (n *noopPushNotificationService) Send(_ context.Context, _ *model.PushSubscription, _ any) error {
	return nil
}

func (n *noopPushNotificationService) PublishToBuyer(_ context.Context, _ int, _ any) error {
	return nil
}
