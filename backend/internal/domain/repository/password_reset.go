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

type PasswordResetRepository interface {
	Create(ctx context.Context, buyerID int, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	Delete(ctx context.Context, tokenHash string) error
	DeleteByBuyerID(ctx context.Context, buyerID int) error
}
