package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type passwordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) repository.PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(ctx context.Context, userID int, role string, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO password_reset_tokens (user_id, user_role, token_hash, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, userID, role, tokenHash, expiresAt)
	return err
}

func (r *passwordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (int, string, time.Time, error) {
	var userID int
	var role string
	var expiresAt time.Time
	query := `SELECT user_id, user_role, expires_at FROM password_reset_tokens WHERE token_hash = $1`
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&userID, &role, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", time.Time{}, nil
		}
		return 0, "", time.Time{}, err
	}
	return userID, role, expiresAt, nil
}

func (r *passwordResetRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM password_reset_tokens WHERE token_hash = $1`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *passwordResetRepository) DeleteAllByUserID(ctx context.Context, userID int, role string) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1 AND user_role = $2`
	_, err := r.db.ExecContext(ctx, query, userID, role)
	return err
}
