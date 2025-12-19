package buyer

import (
	"context"
	"time"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	MaxFailedAttempts = 5
	LockDuration      = 30 * time.Minute
)

// LoginBuyerUseCase defines the interface for buyer login
type LoginBuyerUseCase interface {
	Execute(ctx context.Context, email, password string) (*model.Buyer, error)
}

// loginBuyerUseCase handles buyer login
type loginBuyerUseCase struct {
	buyerRepo repository.BuyerRepository
	authRepo  repository.AuthenticationRepository
}

// NewLoginBuyerUseCase creates a new instance of LoginBuyerUseCase
func NewLoginBuyerUseCase(buyerRepo repository.BuyerRepository, authRepo repository.AuthenticationRepository) LoginBuyerUseCase {
	return &loginBuyerUseCase{
		buyerRepo: buyerRepo,
		authRepo:  authRepo,
	}
}

// Execute authenticates a buyer
func (uc *loginBuyerUseCase) Execute(ctx context.Context, email, password string) (*model.Buyer, error) {
	// Find authentication by email
	auth, err := uc.authRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, &apperrors.UnauthorizedError{Message: "invalid credentials"}
	}

	// Check if account is locked
	if auth.LockedUntil != nil && time.Now().Before(*auth.LockedUntil) {
		return nil, &apperrors.UnauthorizedError{Message: "account is locked due to too many failed attempts"}
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(password)); err != nil {
		// Increment failed attempts
		_ = uc.authRepo.IncrementFailedAttempts(ctx, auth.ID)

		// Lock account if too many failed attempts
		if auth.FailedAttempts+1 >= MaxFailedAttempts {
			lockUntil := time.Now().Add(LockDuration)
			_ = uc.authRepo.LockAccount(ctx, auth.ID, lockUntil)
			return nil, &apperrors.UnauthorizedError{Message: "account locked due to too many failed attempts"}
		}

		return nil, &apperrors.UnauthorizedError{Message: "invalid credentials"}
	}

	// Update last login and reset failed attempts
	_ = uc.authRepo.UpdateLoginSuccess(ctx, auth.ID, time.Now())

	// Get buyer details
	buyer, err := uc.buyerRepo.FindByID(ctx, auth.BuyerID)
	if err != nil {
		return nil, err
	}

	return buyer, nil
}
