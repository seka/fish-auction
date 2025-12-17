package repository

import (
	"context"
	"time"
)

// PasswordResetRepository defines the interface for password reset token persistence
type PasswordResetRepository interface {
	Create(ctx context.Context, userID int, role string, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(ctx context.Context, tokenHash string) (int, string, time.Time, error) // Returns userID, role, expiresAt
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteAllByUserID(ctx context.Context, userID int, role string) error
}
