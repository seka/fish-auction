package postgres

import (
	"context"
	"database/sql"
	"fmt"

	// Added import
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type pushRepository struct {
	db *sql.DB
}

// NewPushRepository creates a new instance of PushRepository
func NewPushRepository(db *sql.DB) repository.PushRepository {
	return &pushRepository{
		db: db,
	}
}

// getDB returns the transaction if one exists in context, otherwise returns the default DB
func (r *pushRepository) getDB(ctx context.Context) dbExecutor {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return r.db
}

func (r *pushRepository) SaveSubscription(ctx context.Context, sub *model.PushSubscription) error {
	db := r.getDB(ctx)

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

	err := db.QueryRowContext(ctx, query,
		sub.BuyerID, sub.Endpoint, sub.P256dh, sub.Auth,
	).Scan(&sub.ID, &sub.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to save subscription: %w", err)
	}
	return nil
}

func (r *pushRepository) GetSubscriptionsByBuyerID(ctx context.Context, buyerID int) ([]model.PushSubscription, error) {
	db := r.getDB(ctx)
	query := `
		SELECT id, buyer_id, endpoint, p256dh, auth, created_at
		FROM push_subscriptions
		WHERE buyer_id = $1
	`

	rows, err := db.QueryContext(ctx, query, buyerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []model.PushSubscription
	for rows.Next() {
		var sub model.PushSubscription
		if err := rows.Scan(
			&sub.ID, &sub.BuyerID, &sub.Endpoint, &sub.P256dh, &sub.Auth, &sub.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (r *pushRepository) DeleteSubscription(ctx context.Context, endpoint string) error {
	db := r.getDB(ctx)
	query := "DELETE FROM push_subscriptions WHERE endpoint = $1"

	// If endpoint URL is long, Postgres handles text type fine.
	// But sometimes endpoint might differ slightly? No, usually exact match.
	// However, depending on browser, endpoint might be encoded. Assuming exact match for now.

	_, err := db.ExecContext(ctx, query, endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}
