package service

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// PushNotificationService defines the interface for delivering push notifications
type PushNotificationService interface {
	Send(ctx context.Context, sub *model.PushSubscription, payload any) error
}
