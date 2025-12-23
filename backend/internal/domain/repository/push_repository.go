package repository

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// PushRepository defines the interface for managing push subscriptions
type PushRepository interface {
	// SaveSubscription saves or updates a subscription for a buyer
	SaveSubscription(ctx context.Context, sub *model.PushSubscription) error
	// GetSubscriptionsByBuyerID retrieves all subscriptions for a buyer
	GetSubscriptionsByBuyerID(ctx context.Context, buyerID int) ([]model.PushSubscription, error)
	// DeleteSubscription deletes a subscription by endpoint (used when expired/unsubscribed)
	DeleteSubscription(ctx context.Context, endpoint string) error
}
