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

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, token, newPassword string) error
}

type resetPasswordUseCase struct {
	pwdResetRepo repository.PasswordResetRepository
	adminRepo    repository.AdminRepository
}

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

	adminID, role, expiresAt, err := u.pwdResetRepo.FindByTokenHash(ctx, tokenHash) // Changed to tokenHash from hashedToken
	if err != nil {
		return err
	}
	if adminID == 0 || role != "admin" {
		return &errors.UnauthorizedError{Message: "Invalid or expired token"}
	}

	if time.Now().After(expiresAt) {
		_ = u.pwdResetRepo.DeleteByTokenHash(ctx, tokenHash)
		return fmt.Errorf("token expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := u.adminRepo.UpdatePassword(ctx, adminID, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := u.pwdResetRepo.DeleteAllByUserID(ctx, adminID, "admin"); err != nil {
		// Log error
	}

	return nil
}
