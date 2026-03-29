package buyer

import (
	"context"
	"fmt"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// UpdatePasswordUseCase defines the interface for updating a buyer password.
type UpdatePasswordUseCase interface {
	Execute(ctx context.Context, buyerID int, currentPassword, newPassword string) error
}

type updatePasswordUseCase struct {
	authRepo    repository.AuthenticationRepository
	sessionRepo repository.SessionRepository
}

// NewUpdatePasswordUseCase creates a new instance of UpdatePasswordUseCase.
func NewUpdatePasswordUseCase(authRepo repository.AuthenticationRepository, sessionRepo repository.SessionRepository) UpdatePasswordUseCase {
	return &updatePasswordUseCase{
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
}

// Execute updates the buyer password after verifying the current one.
func (uc *updatePasswordUseCase) Execute(ctx context.Context, buyerID int, currentPassword, newPassword string) error {
	auth, err := uc.authRepo.FindByBuyerID(ctx, buyerID)
	if err != nil {
		return fmt.Errorf("failed to find authentication: %w", err)
	}
	if auth == nil {
		return &apperrors.NotFoundError{Resource: "authentication", ID: buyerID}
	}

	// 0. Verify current password
	hp := model.NewHashedPassword(auth.PasswordHash)
	if err := hp.Verify(currentPassword); err != nil {
		return err
	}

	// 1. Validate and hash new password
	newPwd, err := model.NewPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to validate new password format: %w", err)
	}

	hp, err = newPwd.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 2. Update password
	if err := uc.authRepo.UpdatePassword(ctx, buyerID, hp.Raw()); err != nil {
		return fmt.Errorf("failed to update password in repository: %w", err)
	}

	// 3. Invalidate all sessions after password change for security
	if err := uc.sessionRepo.DeleteAllByUserID(ctx, buyerID, model.SessionRoleBuyer); err != nil {
		return fmt.Errorf("failed to invalidate sessions: %w", err)
	}

	return nil
}
