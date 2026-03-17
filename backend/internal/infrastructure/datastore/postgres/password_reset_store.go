package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres/errors"
)

var _ repository.PasswordResetRepository = (*passwordResetStore)(nil)

type passwordResetStore struct {
	db datastore.Database
}

// NewPasswordResetStore creates a new instance of PasswordResetRepository.
func NewPasswordResetStore(db datastore.Database) *passwordResetStore {
	return &passwordResetStore{db: db}
}

func (r *passwordResetStore) Create(ctx context.Context, userID int, role string, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO password_reset_tokens (user_id, user_role, token_hash, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Execute(ctx, query, userID, role, tokenHash, expiresAt)
	if err != nil {
		return dserrors.HandleError(err, "PasswordReset", userID, "Create")
	}
	return nil
}

func (r *passwordResetStore) FindByTokenHash(ctx context.Context, tokenHash string) (int, string, time.Time, error) {
	var userID int
	var role string
	var expiresAt time.Time
	query := `SELECT user_id, user_role, expires_at FROM password_reset_tokens WHERE token_hash = $1`
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&userID, &role, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", time.Time{}, nil
		}
		return 0, "", time.Time{}, dserrors.HandleError(err, "PasswordReset", 0, "FindByTokenHash")
	}
	return userID, role, expiresAt, nil
}

func (r *passwordResetStore) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM password_reset_tokens WHERE token_hash = $1`
	_, err := r.db.Execute(ctx, query, tokenHash)
	if err != nil {
		return dserrors.HandleError(err, "PasswordReset", 0, "DeleteByTokenHash")
	}
	return nil
}

func (r *passwordResetStore) DeleteAllByUserID(ctx context.Context, userID int, role string) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1 AND user_role = $2`
	_, err := r.db.Execute(ctx, query, userID, role)
	if err != nil {
		return dserrors.HandleError(err, "PasswordReset", userID, "DeleteAllByUserID")
	}
	return nil
}
