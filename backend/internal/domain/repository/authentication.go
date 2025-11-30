package repository

import (
	"context"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type AuthenticationRepository interface {
	Create(ctx context.Context, auth *model.Authentication) (*model.Authentication, error)
	FindByEmail(ctx context.Context, email string) (*model.Authentication, error)
	FindByBuyerID(ctx context.Context, buyerID int) (*model.Authentication, error)
	UpdateLoginSuccess(ctx context.Context, id int, loginAt time.Time) error
	IncrementFailedAttempts(ctx context.Context, id int) error
	ResetFailedAttempts(ctx context.Context, id int) error
	LockAccount(ctx context.Context, id int, until time.Time) error
}
