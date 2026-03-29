package auth

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// LoginUseCase defines the interface for user authentication.
type LoginUseCase interface {
	// Execute authenticates a user with the provided email and password.
	Execute(ctx context.Context, email, password string) (*model.Admin, error)
}

// LoginUseCase handles user authentication
type loginUseCase struct {
	adminRepo repository.AdminRepository
}

var _ LoginUseCase = (*loginUseCase)(nil)

// NewLoginUseCase creates a new instance of LoginUseCase
func NewLoginUseCase(adminRepo repository.AdminRepository) LoginUseCase {
	return &loginUseCase{adminRepo: adminRepo}
}

// Execute authenticates a user with the provided email and password.
func (u *loginUseCase) Execute(ctx context.Context, email, password string) (*model.Admin, error) {
	admin, err := u.adminRepo.FindOneByEmail(ctx, email)
	if err != nil {
		var nfErr *apperrors.NotFoundError
		if errors.As(err, &nfErr) {
			return nil, &apperrors.UnauthorizedError{Message: "Invalid credentials"}
		}
		return nil, fmt.Errorf("failed to find admin during login: %w", err)
	}
	if admin == nil {
		return nil, &apperrors.UnauthorizedError{Message: "Invalid credentials"}
	}

	// For login, we only need to verify the password against the stored hash.
	// We use HashedPassword which doesn't enforce complexity rules to avoid locking out existing users.
	hp := model.NewHashedPassword(admin.PasswordHash)
	if err := hp.Verify(password); err != nil {
		return nil, err // Verify already returns UnauthorizedError for mismatches
	}

	return admin, nil
}
