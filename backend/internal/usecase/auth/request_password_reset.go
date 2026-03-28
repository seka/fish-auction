package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

var randRead = rand.Read

// RequestPasswordResetUseCase defines the interface for requesting a password reset.
type RequestPasswordResetUseCase interface {
	// Execute initiates the password reset process for the given email.
	Execute(ctx context.Context, email string) error
}

type requestPasswordResetUseCase struct {
	buyerRepo    repository.BuyerRepository
	pwdResetRepo repository.PasswordResetRepository
	emailService service.BuyerEmailService
	frontendURL  *url.URL
}

var _ RequestPasswordResetUseCase = (*requestPasswordResetUseCase)(nil)

// NewRequestPasswordResetUseCase creates a new instance of RequestPasswordResetUseCase
func NewRequestPasswordResetUseCase(
	buyerRepo repository.BuyerRepository,
	pwdResetRepo repository.PasswordResetRepository,
	emailService service.BuyerEmailService,
	frontendURL *url.URL,
) RequestPasswordResetUseCase {
	return &requestPasswordResetUseCase{
		buyerRepo:    buyerRepo,
		pwdResetRepo: pwdResetRepo,
		emailService: emailService,
		frontendURL:  frontendURL,
	}
}

func (u *requestPasswordResetUseCase) Execute(ctx context.Context, email string) error {
	buyer, err := u.buyerRepo.FindByEmail(ctx, email)
	if err != nil {
		// Security: Don't reveal if user exists.
		// Return nil even if not found.
		return nil
	}
	if buyer == nil {
		return nil
	}

	// 1. Generate secure token
	tokenBytes := make([]byte, 32)
	if _, err := randRead(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// 2. Hash token for DB
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// 3. Save to DB (expires in 30 mins)
	expiresAt := time.Now().Add(30 * time.Minute)
	// Invalidate old tokens for this user first
	if err := u.pwdResetRepo.DeleteAllByUserID(ctx, buyer.ID, "buyer"); err != nil {
		return fmt.Errorf("failed to delete old tokens: %w", err)
	}
	if err = u.pwdResetRepo.Create(ctx, buyer.ID, "buyer", tokenHash, expiresAt); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// 4. Send Email
	resetURL := u.frontendURL.JoinPath("/login/reset_password")
	q := resetURL.Query()
	q.Set("token", token)
	resetURL.RawQuery = q.Encode()

	if err := u.emailService.SendBuyerPasswordReset(ctx, email, resetURL.String()); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
