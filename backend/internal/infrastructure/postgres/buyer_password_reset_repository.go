package postgres

import (
	"context"
	"database/sql"
	"errors" // Added import for errors package
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type buyerPasswordResetRepository struct {
	db *sql.DB
}

func NewBuyerPasswordResetRepository(db *sql.DB) repository.BuyerPasswordResetRepository {
	return &buyerPasswordResetRepository{db: db}
}

func (r *buyerPasswordResetRepository) Create(ctx context.Context, buyerID int, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO buyer_password_reset_tokens (buyer_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, buyerID, tokenHash, expiresAt)
	return err
}

func (r *buyerPasswordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (int, time.Time, error) {
	var buyerID int
	var expiresAt time.Time
	query := `SELECT buyer_id, expires_at FROM buyer_password_reset_tokens WHERE token_hash = $1`
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&buyerID, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, time.Time{}, nil // Return zero values if not found (or custom error)
		}
		return 0, time.Time{}, err
	}
	return buyerID, expiresAt, nil
}

func (r *buyerPasswordResetRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM buyer_password_reset_tokens WHERE token_hash = $1`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *buyerPasswordResetRepository) DeleteAllByBuyerID(ctx context.Context, buyerID int) error {
	query := `DELETE FROM buyer_password_reset_tokens WHERE buyer_id = $1`
	_, err := r.db.ExecContext(ctx, query, buyerID)
	return err
}
