package repository

import (
	"context"
	"time"
)

// AdminPasswordResetRepository defines the interface for password reset token persistence for admins
type AdminPasswordResetRepository interface {
	Create(ctx context.Context, adminID int, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(ctx context.Context, tokenHash string) (int, time.Time, error) // Returns adminID, expiresAt
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteAllByAdminID(ctx context.Context, adminID int) error
}
