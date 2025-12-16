package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type adminPasswordResetRepository struct {
	db *sql.DB
}

func NewAdminPasswordResetRepository(db *sql.DB) repository.AdminPasswordResetRepository {
	return &adminPasswordResetRepository{db: db}
}

func (r *adminPasswordResetRepository) Create(ctx context.Context, adminID int, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO admin_password_reset_tokens (admin_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, adminID, tokenHash, expiresAt)
	return err
}

func (r *adminPasswordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (int, time.Time, error) {
	var adminID int
	var expiresAt time.Time
	query := `SELECT admin_id, expires_at FROM admin_password_reset_tokens WHERE token_hash = $1`
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&adminID, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, time.Time{}, nil
		}
		return 0, time.Time{}, err
	}
	return adminID, expiresAt, nil
}

func (r *adminPasswordResetRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM admin_password_reset_tokens WHERE token_hash = $1`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *adminPasswordResetRepository) DeleteAllByAdminID(ctx context.Context, adminID int) error {
	query := `DELETE FROM admin_password_reset_tokens WHERE admin_id = $1`
	_, err := r.db.ExecContext(ctx, query, adminID)
	return err
}
