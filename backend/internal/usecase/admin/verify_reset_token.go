package admin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// VerifyResetTokenUseCase defines the interface for verifying an admin reset token.
type VerifyResetTokenUseCase interface {
	Execute(ctx context.Context, token string) error
}

type verifyResetTokenUseCase struct {
	pwdResetRepo repository.PasswordResetRepository
}

var _ VerifyResetTokenUseCase = (*verifyResetTokenUseCase)(nil)

// NewVerifyResetTokenUseCase creates a new VerifyResetTokenUseCase instance.
func NewVerifyResetTokenUseCase(pwdResetRepo repository.PasswordResetRepository) VerifyResetTokenUseCase {
	return &verifyResetTokenUseCase{pwdResetRepo: pwdResetRepo}
}

func (u *verifyResetTokenUseCase) Execute(ctx context.Context, token string) error {
	// 1. Hash token
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// 2. Find token in DB
	resetToken, err := u.pwdResetRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return err
	}
	if resetToken == nil || resetToken.Role != "admin" {
		return &errors.UnauthorizedError{Message: "Invalid or expired token"}
	}

	// 3. Check expiry
	if time.Now().After(resetToken.ExpiresAt) {
		_ = u.pwdResetRepo.DeleteByTokenHash(ctx, tokenHash)
		return fmt.Errorf("token expired")
	}

	return nil
}
