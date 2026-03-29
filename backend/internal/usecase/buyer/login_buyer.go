package buyer

import (
	"context"
	"fmt"
	"time"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase"
)

// LockDuration and MaxFailedAttempts are now defined in central usecase package.

// LoginBuyerUseCase defines the interface for buyer login
type LoginBuyerUseCase interface {
	Execute(ctx context.Context, email, password string) (*model.Buyer, error)
}

// LoginBuyerUseCase handles buyer login
type loginBuyerUseCase struct {
	buyerRepo repository.BuyerRepository
	authRepo  repository.AuthenticationRepository
}

var _ LoginBuyerUseCase = (*loginBuyerUseCase)(nil)

// NewLoginBuyerUseCase creates a new instance of LoginBuyerUseCase
func NewLoginBuyerUseCase(buyerRepo repository.BuyerRepository, authRepo repository.AuthenticationRepository) LoginBuyerUseCase {
	return &loginBuyerUseCase{buyerRepo: buyerRepo, authRepo: authRepo}
}

// Execute authenticates a buyer
func (uc *loginBuyerUseCase) Execute(ctx context.Context, email, password string) (*model.Buyer, error) {
	// 0. Validate password format (optional for login, but good for consistency)
	pwd, err := model.NewPassword(password)
	if err != nil {
		return nil, &apperrors.UnauthorizedError{Message: "invalid credentials"}
	}

	// Find authentication by email
	auth, err := uc.authRepo.FindByEmail(ctx, email)
	if err != nil {
		// We return UnauthorizedError to mask whether the user exists or not.
		// However, we can log the real error or wrap it if needed for internal tracing.
		return nil, &apperrors.UnauthorizedError{Message: "invalid credentials"}
	}
	if auth == nil {
		return nil, &apperrors.UnauthorizedError{Message: "invalid credentials"}
	}

	// Check if account is locked
	if auth.LockedUntil != nil && time.Now().Before(*auth.LockedUntil) {
		return nil, &apperrors.UnauthorizedError{Message: "account is locked due to too many failed attempts"}
	}

	// Verify password using Password VO
	if err := pwd.CompareWithHash(auth.PasswordHash); err != nil {
		// Increment failed attempts
		_ = uc.authRepo.IncrementFailedAttempts(ctx, auth.ID)

		// Lock account if too many failed attempts
		if auth.FailedAttempts+1 >= usecase.MaxFailedLoginAttempts {
			lockUntil := time.Now().Add(usecase.AccountLockDuration)
			_ = uc.authRepo.LockAccount(ctx, auth.ID, lockUntil)
			return nil, &apperrors.UnauthorizedError{Message: "account locked due to too many failed attempts"}
		}

		return nil, &apperrors.UnauthorizedError{Message: "invalid credentials"}
	}

	// Update last login and reset failed attempts
	if err := uc.authRepo.UpdateLoginSuccess(ctx, auth.ID, time.Now()); err != nil {
		return nil, fmt.Errorf("failed to update login success: %w", err)
	}

	// Get buyer details
	buyer, err := uc.buyerRepo.FindByID(ctx, auth.BuyerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find buyer details: %w", err)
	}

	return buyer, nil
}
