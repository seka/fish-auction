package auth

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

// ResetPasswordUseCase defines the interface for resetting a password.
type ResetPasswordUseCase interface {
	// Execute resets the password using the reset token.
	Execute(ctx context.Context, token, newPassword string) error
}

type resetPasswordUseCase struct {
	pwdResetRepo repository.PasswordResetRepository
	authRepo     repository.AuthenticationRepository
}

var _ ResetPasswordUseCase = (*resetPasswordUseCase)(nil)

// NewResetPasswordUseCase creates a new instance of ResetPasswordUseCase
func NewResetPasswordUseCase(
	pwdResetRepo repository.PasswordResetRepository,
	authRepo repository.AuthenticationRepository,
) ResetPasswordUseCase {
	return &resetPasswordUseCase{
		pwdResetRepo: pwdResetRepo,
		authRepo:     authRepo,
	}
}

func (u *resetPasswordUseCase) Execute(ctx context.Context, token, newPassword string) error {
	// 1. Hash token to verify
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// 2. Find token in DB
	resetToken, err := u.pwdResetRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return err
	}
	if resetToken == nil || resetToken.Role != "buyer" { // Check role
		return &errors.UnauthorizedError{Message: "Invalid or expired token"}
	}

	// 3. Check expiry
	if time.Now().After(resetToken.ExpiresAt) {
		// Clean up expired token
		_ = u.pwdResetRepo.DeleteByTokenHash(ctx, tokenHash)
		return &errors.UnauthorizedError{Message: "Invalid or expired token"}
	}

	// 4. Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 5. Update user password
	if err := u.authRepo.UpdatePassword(ctx, resetToken.UserID, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate token
	_ = u.pwdResetRepo.DeleteAllByUserID(ctx, resetToken.UserID, "buyer")

	return nil
}
