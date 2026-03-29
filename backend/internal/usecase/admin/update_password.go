package admin

import (
	"context"
	"fmt"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// UpdatePasswordUseCase defines the interface for updating an admin password.
type UpdatePasswordUseCase interface {
	// Execute updates the admin password for the given ID.
	Execute(ctx context.Context, id int, currentPassword, newPassword string) error
}

var _ UpdatePasswordUseCase = (*updatePasswordUseCase)(nil)

type updatePasswordUseCase struct {
	adminRepo   repository.AdminRepository
	sessionRepo repository.SessionRepository
}

// NewUpdatePasswordUseCase creates a new instance of UpdatePasswordUseCase.
func NewUpdatePasswordUseCase(adminRepo repository.AdminRepository, sessionRepo repository.SessionRepository) UpdatePasswordUseCase {
	return &updatePasswordUseCase{
		adminRepo:   adminRepo,
		sessionRepo: sessionRepo,
	}
}

// Execute updates the admin password after verifying the current one.
func (uc *updatePasswordUseCase) Execute(ctx context.Context, id int, currentPassword, newPassword string) error {
	admin, err := uc.adminRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find admin: %w", err)
	}
	if admin == nil {
		return &apperrors.NotFoundError{Resource: "admin", ID: id}
	}

	// 0. Verify current password
	// We use HashedPassword to allow existing simple passwords to be verified.
	hp := model.NewHashedPassword(admin.PasswordHash)
	if err := hp.Verify(currentPassword); err != nil {
		// Replace standard UnauthorizedError with more specific context if needed,
		// but Verify already returns a proper UnauthorizedError.
		return err
	}

	// 1. Validate and hash new password
	// NewPassword enforces strict complexity requirements.
	newPwd, err := model.NewPassword(newPassword)
	if err != nil {
		return err // Returns ValidationError if fails
	}

	hashedPwd, err := newPwd.Hash()
	if err != nil {
		return err
	}

	// 2. Update password
	if err := uc.adminRepo.UpdatePassword(ctx, id, hashedPwd.Raw()); err != nil {
		return fmt.Errorf("failed to update password in repository: %w", err)
	}

	// 3. Invalidate all sessions after password change for security
	if err := uc.sessionRepo.DeleteAllByUserID(ctx, id, model.SessionRoleAdmin); err != nil {
		return fmt.Errorf("failed to invalidate sessions after password update: %w", err)
	}

	return nil
}
