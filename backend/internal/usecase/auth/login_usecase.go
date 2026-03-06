package auth

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// LoginUseCase defines the interface for user authentication
type LoginUseCase interface {
	Execute(ctx context.Context, email, password string) (*model.Admin, error)
}

// loginUseCase handles user authentication
type loginUseCase struct {
	adminRepo repository.AdminRepository
}

// NewLoginUseCase creates a new instance of LoginUseCase
func NewLoginUseCase(adminRepo repository.AdminRepository) LoginUseCase {
	return &loginUseCase{adminRepo: adminRepo}
}

// Execute authenticates a user with the provided password
func (u *loginUseCase) Execute(ctx context.Context, email, password string) (*model.Admin, error) {
	admin, err := u.adminRepo.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, nil // Not found
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil // Invalid password
	}

	return admin, nil
}
