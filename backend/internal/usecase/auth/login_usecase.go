package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"time"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

const (
	// MaxAdminFailedLoginAttempts is the number of consecutive failed login
	// attempts before an admin account is locked.
	MaxAdminFailedLoginAttempts = 5

	// AdminAccountLockDuration is the duration for which an admin account is locked
	// after exceeding MaxAdminFailedLoginAttempts.
	AdminAccountLockDuration = 30 * time.Minute
)

// LoginUseCase defines the interface for user authentication.
type LoginUseCase interface {
	// Execute authenticates a user with the provided email and password.
	Execute(ctx context.Context, email, password string) (*model.Admin, error)
}

// loginUseCase handles admin authentication with lockout support.
type loginUseCase struct {
	adminRepo repository.AdminRepository
	clock     service.Clock
}

var _ LoginUseCase = (*loginUseCase)(nil)

// NewLoginUseCase creates a new instance of LoginUseCase
func NewLoginUseCase(adminRepo repository.AdminRepository, clock service.Clock) LoginUseCase {
	return &loginUseCase{adminRepo: adminRepo, clock: clock}
}

// Execute authenticates an admin with the provided email and password.
func (u *loginUseCase) Execute(ctx context.Context, email, password string) (*model.Admin, error) {
	admin, err := u.adminRepo.FindOneByEmail(ctx, email)
	if err != nil {
		var nfErr *apperrors.NotFoundError
		if errors.As(err, &nfErr) {
			slog.WarnContext(ctx, "auth: admin login failed", "reason", "user_not_found", "email_hash", emailHash(email))
			return nil, &apperrors.UnauthorizedError{Message: "Invalid credentials"}
		}
		return nil, fmt.Errorf("failed to find admin during login: %w", err)
	}
	if admin == nil {
		slog.WarnContext(ctx, "auth: admin login failed", "reason", "user_not_found", "email_hash", emailHash(email))
		return nil, &apperrors.UnauthorizedError{Message: "Invalid credentials"}
	}

	now := u.clock.Now()
	if admin.LockedUntil != nil && now.Before(*admin.LockedUntil) {
		slog.WarnContext(ctx, "auth: admin login failed", "reason", "account_locked", "email_hash", emailHash(email))
		return nil, &apperrors.UnauthorizedError{Message: "account is locked due to too many failed attempts"}
	}

	// For login, we only need to verify the password against the stored hash.
	// We use HashedPassword which doesn't enforce complexity rules to avoid locking out existing users.
	hp := model.NewHashedPassword(admin.PasswordHash)
	if err := hp.Verify(password); err != nil {
		_ = u.adminRepo.IncrementFailedAttempts(ctx, admin.ID)

		newAttempts := admin.FailedAttempts + 1
		slog.WarnContext(ctx, "auth: admin login failed", "reason", "bad_password", "email_hash", emailHash(email), "attempts", newAttempts)

		if newAttempts >= MaxAdminFailedLoginAttempts {
			lockUntil := u.clock.Now().Add(AdminAccountLockDuration)
			_ = u.adminRepo.LockAccount(ctx, admin.ID, lockUntil)
			return nil, &apperrors.UnauthorizedError{Message: "account locked due to too many failed attempts"}
		}

		return nil, err
	}

	if err := u.adminRepo.UpdateLoginSuccess(ctx, admin.ID, u.clock.Now()); err != nil {
		return nil, fmt.Errorf("failed to update login success: %w", err)
	}

	return admin, nil
}

// emailHash returns a 4-byte hex prefix of the SHA-256 hash of the email
// for privacy-preserving log correlation.
func emailHash(email string) string {
	h := sha256.Sum256([]byte(email))
	return fmt.Sprintf("%x", h[:4])
}
