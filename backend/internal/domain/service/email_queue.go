package service

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// BuyerEmailQueue defines the interface for enqueueing buyer-related emails
type BuyerEmailQueue interface {
	EnqueueBuyerPasswordReset(ctx context.Context, to, url string) error
	Dequeue(ctx context.Context, waitTimeSeconds int32) ([]*model.JobMessage, error)
	DeleteMessage(ctx context.Context, message *model.JobMessage) error
}

// AdminEmailQueue defines the interface for enqueueing admin-related emails
type AdminEmailQueue interface {
	EnqueueAdminPasswordReset(ctx context.Context, to, url string) error
	Dequeue(ctx context.Context, waitTimeSeconds int32) ([]*model.JobMessage, error)
	DeleteMessage(ctx context.Context, message *model.JobMessage) error
}
