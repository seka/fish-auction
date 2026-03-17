package postgres

import (
	"context"
	"fmt"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres/errors"
)

// PushStore implements repository.PushRepository.
type PushStore struct {
	db datastore.Database
}

var _ repository.PushRepository = (*PushStore)(nil)

// NewPushStore creates a new PushStore instance.
func NewPushStore(db datastore.Database) *PushStore {
	return &PushStore{db: db}
}

// SaveSubscription stores or updates a push subscription.
func (r *PushStore) SaveSubscription(ctx context.Context, sub *model.PushSubscription) error {

	// Upsert subscription based on endpoint
	query := `
		INSERT INTO push_subscriptions (buyer_id, endpoint, p256dh, auth)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (endpoint) DO UPDATE
		SET buyer_id = EXCLUDED.buyer_id,
		    p256dh = EXCLUDED.p256dh,
		    auth = EXCLUDED.auth,
		    created_at = CURRENT_TIMESTAMP
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		sub.BuyerID, sub.Endpoint, sub.P256dh, sub.Auth,
	).Scan(&sub.ID, &sub.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to save subscription: %w", dserrors.HandleError(err, "PushSubscription", 0, "SaveSubscription"))
	}
	return nil
}

// GetSubscriptionsByBuyerID returns all push subscriptions for a buyer.
func (r *PushStore) GetSubscriptionsByBuyerID(ctx context.Context, buyerID int) ([]model.PushSubscription, error) {
	query := `
		SELECT id, buyer_id, endpoint, p256dh, auth, created_at
		FROM push_subscriptions
		WHERE buyer_id = $1
	`

	rows, err := r.db.Query(ctx, query, buyerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", dserrors.HandleError(err, "PushSubscription", buyerID, "GetSubscriptionsByBuyerID"))
	}
	defer func() { _ = rows.Close() }()

	var subs []model.PushSubscription
	for rows.Next() {
		var sub model.PushSubscription
		if err := rows.Scan(
			&sub.ID, &sub.BuyerID, &sub.Endpoint, &sub.P256dh, &sub.Auth, &sub.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", dserrors.HandleError(err, "PushSubscription", buyerID, "GetSubscriptionsByBuyerID"))
		}
		subs = append(subs, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, dserrors.HandleError(err, "PushSubscription", buyerID, "GetSubscriptionsByBuyerID")
	}
	return subs, nil
}

// DeleteSubscription removes a push subscription by its endpoint.
func (r *PushStore) DeleteSubscription(ctx context.Context, endpoint string) error {
	query := "DELETE FROM push_subscriptions WHERE endpoint = $1"

	// If endpoint URL is long, Postgres handles text type fine.
	// But sometimes endpoint might differ slightly? No, usually exact match.
	// However, depending on browser, endpoint might be encoded. Assuming exact match for now.

	_, err := r.db.Execute(ctx, query, endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", dserrors.HandleError(err, "PushSubscription", 0, "DeleteSubscription"))
	}
	return nil
}
