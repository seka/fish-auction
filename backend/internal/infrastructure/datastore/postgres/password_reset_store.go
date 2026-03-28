package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres/errors"
)

var _ repository.PasswordResetRepository = (*PasswordResetStore)(nil)

// PasswordResetStore implements repository.PasswordResetRepository using PostgreSQL.
type PasswordResetStore struct {
	db datastore.Database
}

// NewPasswordResetStore creates a new instance of PasswordResetRepository
func NewPasswordResetStore(db datastore.Database) *PasswordResetStore {
	return &PasswordResetStore{db: db}
}

// Create stores a new password reset token.
func (r *PasswordResetStore) Create(ctx context.Context, userID int, role, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO password_reset_tokens (user_id, user_role, token_hash, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Execute(ctx, query, userID, role, tokenHash, expiresAt)
	if err != nil {
		return dserrors.HandleError(err, "PasswordReset", userID, "Create")
	}
	return nil
}

// FindByTokenHash returns password reset information by token hash.
func (r *PasswordResetStore) FindByTokenHash(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error) {
	query := `SELECT user_id, user_role, expires_at FROM password_reset_tokens WHERE token_hash = $1`
	var res model.PasswordResetToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&res.UserID, &res.Role, &res.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, dserrors.HandleError(err, "PasswordReset", 0, "FindByTokenHash")
	}
	res.TokenHash = tokenHash
	return &res, nil
}

// DeleteByTokenHash removes a password reset token.
func (r *PasswordResetStore) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM password_reset_tokens WHERE token_hash = $1`
	_, err := r.db.Execute(ctx, query, tokenHash)
	if err != nil {
		return dserrors.HandleError(err, "PasswordReset", 0, "DeleteByTokenHash")
	}
	return nil
}

// DeleteAllByUserID removes all password reset tokens for a user.
func (r *PasswordResetStore) DeleteAllByUserID(ctx context.Context, userID int, role string) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1 AND user_role = $2`
	_, err := r.db.Execute(ctx, query, userID, role)
	if err != nil {
		return dserrors.HandleError(err, "PasswordReset", userID, "DeleteAllByUserID")
	}
	return nil
}
