package service

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// PushNotificationQueue defines the interface for enqueueing push notifications
type PushNotificationQueue interface {
	Enqueue(ctx context.Context, buyerID int, payload any) error
	Dequeue(ctx context.Context, waitTimeSeconds int32) ([]*model.JobMessage, error)
	DeleteMessage(ctx context.Context, message *model.JobMessage) error
}
