package push_notification

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type noopPushNotificationService struct{}

func (n *noopPushNotificationService) Send(_ context.Context, _ *model.PushSubscription, _ any) error {
	return nil
}

