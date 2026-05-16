package repository

import (
	"context"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// AdminRepository defines the interface for admin data persistence
type AdminRepository interface {
	FindOneByEmail(ctx context.Context, email string) (*model.Admin, error)
	FindByID(ctx context.Context, id int) (*model.Admin, error)
	Create(ctx context.Context, admin *model.Admin) error
	Count(ctx context.Context) (int, error)
	UpdatePassword(ctx context.Context, id int, passwordHash string) error
	IncrementFailedAttempts(ctx context.Context, id int) error
	LockAccount(ctx context.Context, id int, until time.Time) error
	UpdateLoginSuccess(ctx context.Context, id int, now time.Time) error
}
