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
	// Find authentication by email
	auth, err := uc.authRepo.FindByEmail(ctx, email)
	if err != nil {
		// Return internal error for tracing if it's a real DB fault
		return nil, fmt.Errorf("failed to find authentication during login: %w", err)
	}
	if auth == nil {
		// Only mask user existence by returning Unauthorized
		return nil, &apperrors.UnauthorizedError{Message: "invalid credentials"}
	}

	// Check if account is locked
	if auth.LockedUntil != nil && time.Now().Before(*auth.LockedUntil) {
		return nil, &apperrors.UnauthorizedError{Message: "account is locked due to too many failed attempts"}
	}

	// Verify password using HashedPassword to allow existing simple passwords
	hp := model.NewHashedPassword(auth.PasswordHash)
	if err := hp.Verify(password); err != nil {
		// Increment failed attempts
		_ = uc.authRepo.IncrementFailedAttempts(ctx, auth.ID)

		// Lock account if too many failed attempts
		if auth.FailedAttempts+1 >= usecase.MaxFailedLoginAttempts {
			lockUntil := time.Now().Add(usecase.AccountLockDuration)
			_ = uc.authRepo.LockAccount(ctx, auth.ID, lockUntil)
			return nil, &apperrors.UnauthorizedError{Message: "account locked due to too many failed attempts"}
		}

		return nil, err // Verify already returns UnauthorizedError for mismatches
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
