package admin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
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
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	resetToken, err := u.pwdResetRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return err
	}
	if resetToken == nil || resetToken.Role != "admin" {
		return &errors.UnauthorizedError{Message: "Invalid or expired token"}
	}

	if time.Now().After(resetToken.ExpiresAt) {
		_ = u.pwdResetRepo.DeleteByTokenHash(ctx, tokenHash)
		return fmt.Errorf("token expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := u.adminRepo.UpdatePassword(ctx, resetToken.UserID, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate token
	_ = u.pwdResetRepo.DeleteAllByUserID(ctx, resetToken.UserID, "admin")

	return nil
}
