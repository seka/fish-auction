package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type passwordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) repository.PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(ctx context.Context, buyerID int, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (buyer_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, buyerID, tokenHash, expiresAt)
	return err
}

func (r *passwordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*repository.PasswordResetToken, error) {
	query := `
		SELECT id, buyer_id, token_hash, expires_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1
	`
	row := r.db.QueryRowContext(ctx, query, tokenHash)

	var token repository.PasswordResetToken
	err := row.Scan(
		&token.ID,
		&token.BuyerID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *passwordResetRepository) Delete(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM password_reset_tokens WHERE token_hash = $1`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *passwordResetRepository) DeleteByBuyerID(ctx context.Context, buyerID int) error {
	query := `DELETE FROM password_reset_tokens WHERE buyer_id = $1`
	_, err := r.db.ExecContext(ctx, query, buyerID)
	return err
}
