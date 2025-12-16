package repository

import (
	"context"
	"time"
)

type PasswordResetToken struct {
	ID        string
	BuyerID   int
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// BuyerPasswordResetRepository defines the interface for password reset token persistence for buyers
type BuyerPasswordResetRepository interface {
	Create(ctx context.Context, buyerID int, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(ctx context.Context, tokenHash string) (int, time.Time, error) // Returns buyerID, expiresAt
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteAllByBuyerID(ctx context.Context, buyerID int) error
}
