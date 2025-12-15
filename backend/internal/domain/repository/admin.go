package repository

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
)

// AdminRepository defines the interface for admin data persistence
type AdminRepository interface {
	FindOneByEmail(ctx context.Context, email string) (*entity.Admin, error)
	FindByID(ctx context.Context, id int) (*entity.Admin, error)
	Create(ctx context.Context, admin *entity.Admin) error
	Count(ctx context.Context) (int, error)
	UpdatePassword(ctx context.Context, id int, passwordHash string) error
}
