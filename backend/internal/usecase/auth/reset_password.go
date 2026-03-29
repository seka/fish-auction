package auth

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

// ResetPasswordUseCase defines the interface for resetting a password.
type ResetPasswordUseCase interface {
	// Execute resets the password using the reset token.
	Execute(ctx context.Context, token, newPassword string) error
}

type resetPasswordUseCase struct {
	pwdResetRepo repository.PasswordResetRepository
	authRepo     repository.AuthenticationRepository
	txMgr        repository.TransactionManager
}

var _ ResetPasswordUseCase = (*resetPasswordUseCase)(nil)

// NewResetPasswordUseCase creates a new instance of ResetPasswordUseCase
func NewResetPasswordUseCase(
	pwdResetRepo repository.PasswordResetRepository,
	authRepo repository.AuthenticationRepository,
	txMgr repository.TransactionManager,
) ResetPasswordUseCase {
	return &resetPasswordUseCase{
		pwdResetRepo: pwdResetRepo,
		authRepo:     authRepo,
		txMgr:        txMgr,
	}
}

func (u *resetPasswordUseCase) Execute(ctx context.Context, token, newPassword string) error {
	// 0. Validate new password
	newPwd, err := model.NewPassword(newPassword)
	if err != nil {
		return err
	}

	// 1. Hash token to verify
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// 2. Find token in DB
	resetToken, err := u.pwdResetRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to find reset token: %w", err)
	}
	if resetToken == nil || resetToken.Role != "buyer" { // Check role
		return &apperrors.UnauthorizedError{Message: "Invalid or expired token"}
	}

	// 3. Check expiry
	if time.Now().After(resetToken.ExpiresAt) {
		// Clean up expired token
		_ = u.pwdResetRepo.DeleteByTokenHash(ctx, tokenHash)
		return &apperrors.UnauthorizedError{Message: "Invalid or expired token"}
	}

	// 4. Hash new password
	hashedPwd, err := newPwd.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 5. Atomic Update password and Invalidate token
	err = u.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// Update user password
		if err := u.authRepo.UpdatePassword(txCtx, resetToken.UserID, hashedPwd.Raw()); err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}

		// Invalidate token
		if err := u.pwdResetRepo.DeleteAllByUserID(txCtx, resetToken.UserID, "buyer"); err != nil {
			return fmt.Errorf("failed to invalidate reset token after successful reset: %w", err)
		}

		return nil
	})

	return err
}
