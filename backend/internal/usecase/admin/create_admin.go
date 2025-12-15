package admin

import (
	"context"
	"fmt"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// CreateAdminUseCase defines the interface for creating an admin
type CreateAdminUseCase interface {
	Execute(ctx context.Context, email, password string) error
	Count(ctx context.Context) (int, error)
}

type createAdminUseCase struct {
	adminRepo repository.AdminRepository
}

// NewCreateAdminUseCase creates a new instance of CreateAdminUseCase
func NewCreateAdminUseCase(adminRepo repository.AdminRepository) CreateAdminUseCase {
	return &createAdminUseCase{adminRepo: adminRepo}
}

func (uc *createAdminUseCase) Execute(ctx context.Context, email, password string) error {
	// Check if email already exists? (Optional, repo will error on unique constraint, but nicer here)
	existing, err := uc.adminRepo.FindOneByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check existing admin: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("admin already exists with email: %s", email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin := &entity.Admin{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := uc.adminRepo.Create(ctx, admin); err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	return nil
}

func (uc *createAdminUseCase) Count(ctx context.Context) (int, error) {
	return uc.adminRepo.Count(ctx)
}
