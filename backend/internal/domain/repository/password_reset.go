package repository

import (
	"context"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// PasswordResetRepository defines the interface for password reset token persistence
type PasswordResetRepository interface {
	Create(ctx context.Context, userID int, role string, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteAllByUserID(ctx context.Context, userID int, role string) error
}
