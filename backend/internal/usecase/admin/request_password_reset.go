package admin

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

var randRead = rand.Read

// RequestPasswordResetUseCase defines the interface for requesting an admin password reset.
type RequestPasswordResetUseCase interface {
	// Execute initiates the password reset process for the given email.
	Execute(ctx context.Context, email string) error
}

type requestPasswordResetUseCase struct {
	adminRepo    repository.AdminRepository
	pwdResetRepo repository.PasswordResetRepository
	emailService service.AdminEmailService
	frontendURL  *url.URL
}

var _ RequestPasswordResetUseCase = (*requestPasswordResetUseCase)(nil)

// NewRequestPasswordResetUseCase creates a new RequestPasswordResetUseCase instance.
func NewRequestPasswordResetUseCase(
	adminRepo repository.AdminRepository,
	pwdResetRepo repository.PasswordResetRepository,
	emailService service.AdminEmailService,
	frontendURL *url.URL,
) RequestPasswordResetUseCase {
	return &requestPasswordResetUseCase{
		adminRepo:    adminRepo,
		pwdResetRepo: pwdResetRepo,
		emailService: emailService,
		frontendURL:  frontendURL,
	}
}

func (u *requestPasswordResetUseCase) Execute(ctx context.Context, email string) error {
	admin, err := u.adminRepo.FindOneByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to find admin by email: %w", err)
	}
	if admin == nil {
		// Business logic: if admin not found, tell the handler.
		// Obfuscation is handler responsibility.
		return &apperrors.NotFoundError{Resource: "admin", ID: email}
	}

	// 1. Generate secure token
	tokenBytes := make([]byte, 32)
	if _, err := randRead(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate secure token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// 2. Hash token for DB
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// 3. Save to DB (expires in 30 mins)
	expiresAt := time.Now().Add(30 * time.Minute)
	// Invalidate old tokens for this user first
	if err := u.pwdResetRepo.DeleteAllByUserID(ctx, admin.ID, "admin"); err != nil {
		return fmt.Errorf("failed to invalidate old reset tokens: %w", err)
	}
	if err = u.pwdResetRepo.Create(ctx, admin.ID, "admin", tokenHash, expiresAt); err != nil {
		return fmt.Errorf("failed to create new reset token: %w", err)
	}

	// 4. Send Email
	resetURL := u.frontendURL.JoinPath("/login/admin/reset_password")
	q := resetURL.Query()
	q.Set("token", token)
	resetURL.RawQuery = q.Encode()

	if err := u.emailService.SendAdminPasswordReset(ctx, email, resetURL.String()); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}
