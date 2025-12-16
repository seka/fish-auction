package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, token, newPassword string) error
}

type resetPasswordUseCase struct {
	pwdResetRepo repository.PasswordResetRepository
	authRepo     repository.AuthenticationRepository
}

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
	storedToken, err := u.pwdResetRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to find token: %w", err)
	}
	if storedToken == nil {
		return fmt.Errorf("invalid or expired token")
	}

	// 3. Check expiry
	if time.Now().After(storedToken.ExpiresAt) {
		// Clean up expired token
		_ = u.pwdResetRepo.Delete(ctx, tokenHash)
		return fmt.Errorf("token expired")
	}

	// 4. Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 5. Update user password
	if err := u.authRepo.UpdatePassword(ctx, storedToken.BuyerID, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// 6. Invalidate token
	if err := u.pwdResetRepo.Delete(ctx, tokenHash); err != nil {
		// Log error but don't fail, critical part is done
		// log.Printf("failed to delete token: %v", err)
	}

	return nil
}
