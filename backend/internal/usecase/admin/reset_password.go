package admin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ResetPasswordUseCase defines the interface for resetting an admin password.
type ResetPasswordUseCase interface {
	// Execute resets the admin password using the reset token.
	Execute(ctx context.Context, token, newPassword string) error
}

type resetPasswordUseCase struct {
	pwdResetRepo repository.PasswordResetRepository
	adminRepo    repository.AdminRepository
}

var _ ResetPasswordUseCase = (*resetPasswordUseCase)(nil)

// NewResetPasswordUseCase creates a new ResetPasswordUseCase instance.
func NewResetPasswordUseCase(
	pwdResetRepo repository.PasswordResetRepository,
	adminRepo repository.AdminRepository,
) ResetPasswordUseCase {
	return &resetPasswordUseCase{
		pwdResetRepo: pwdResetRepo,
		adminRepo:    adminRepo,
	}
}

func (u *resetPasswordUseCase) Execute(ctx context.Context, token, newPassword string) error {
	// 0. Validate new password
	newPwd, err := model.NewPassword(newPassword)
	if err != nil {
		return err // Returns ValidationError
	}

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	resetToken, err := u.pwdResetRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to find reset token: %w", err)
	}
	if resetToken == nil || resetToken.Role != "admin" {
		return &apperrors.UnauthorizedError{Message: "Invalid or expired token"}
	}

	if time.Now().After(resetToken.ExpiresAt) {
		_ = u.pwdResetRepo.DeleteByTokenHash(ctx, tokenHash)
		return &apperrors.UnauthorizedError{Message: "Invalid or expired token"}
	}

	// 4. Hash new password
	hashedPwd, err := newPwd.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := u.adminRepo.UpdatePassword(ctx, resetToken.UserID, hashedPwd.Raw()); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate token
	if err := u.pwdResetRepo.DeleteAllByUserID(ctx, resetToken.UserID, "admin"); err != nil {
		return fmt.Errorf("failed to invalidate reset token after successful reset: %w", err)
	}

	return nil
}
