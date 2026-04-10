package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

// VerifyResetTokenUseCase defines the interface for verifying a password reset token.
type VerifyResetTokenUseCase interface {
	// Execute verifies the reset token.
	Execute(ctx context.Context, token string) error
}

type verifyResetTokenUseCase struct {
	pwdResetRepo repository.PasswordResetRepository
	clock        service.Clock
}

var _ VerifyResetTokenUseCase = (*verifyResetTokenUseCase)(nil)

// NewVerifyResetTokenUseCase creates a new instance of VerifyResetTokenUseCase.
func NewVerifyResetTokenUseCase(
	pwdResetRepo repository.PasswordResetRepository,
	clock service.Clock,
) VerifyResetTokenUseCase {
	return &verifyResetTokenUseCase{
		pwdResetRepo: pwdResetRepo,
		clock:        clock,
	}
}

func (u *verifyResetTokenUseCase) Execute(ctx context.Context, token string) error {
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
	if u.clock.Now().After(resetToken.ExpiresAt) {
		// Clean up expired token
		_ = u.pwdResetRepo.DeleteByTokenHash(ctx, tokenHash)
		return &errors.UnauthorizedError{Message: "Invalid or expired token"}
	}

	return nil
}
