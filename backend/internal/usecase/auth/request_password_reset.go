package auth

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
	emailMessage "github.com/seka/fish-auction/backend/internal/job/message"
	"github.com/seka/fish-auction/backend/internal/usecase/notification"
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
	publishEmail notification.PublishEmailUseCase
	frontendURL  *url.URL
	txMgr        repository.TransactionManager
	clock        service.Clock
}

var _ RequestPasswordResetUseCase = (*requestPasswordResetUseCase)(nil)

// NewRequestPasswordResetUseCase creates a new instance of RequestPasswordResetUseCase
func NewRequestPasswordResetUseCase(
	buyerRepo repository.BuyerRepository,
	pwdResetRepo repository.PasswordResetRepository,
	publishEmail notification.PublishEmailUseCase,
	frontendURL *url.URL,
	txMgr repository.TransactionManager,
	clock service.Clock,
) RequestPasswordResetUseCase {
	return &requestPasswordResetUseCase{
		buyerRepo:    buyerRepo,
		pwdResetRepo: pwdResetRepo,
		publishEmail: publishEmail,
		frontendURL:  frontendURL,
		txMgr:        txMgr,
		clock:        clock,
	}
}

func (u *requestPasswordResetUseCase) Execute(ctx context.Context, email string) error {
	buyer, err := u.buyerRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to find buyer by email: %w", err)
	}
	if buyer == nil {
		// Business logic: if buyer not found, tell the handler.
		// Obfuscation is handler responsibility.
		return &apperrors.NotFoundError{Resource: "buyer", ID: email}
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

	// 3. Atomic Save to DB (expires in 30 mins)
	expiresAt := u.clock.Now().Add(30 * time.Minute)
	err = u.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// Invalidate old tokens for this user first
		if err := u.pwdResetRepo.DeleteAllByUserID(txCtx, buyer.ID, "buyer"); err != nil {
			return fmt.Errorf("failed to invalidate old reset tokens: %w", err)
		}
		if err = u.pwdResetRepo.Create(txCtx, buyer.ID, "buyer", tokenHash, expiresAt); err != nil {
			return fmt.Errorf("failed to create new reset token: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 4. Enqueue Email Job
	resetURL := u.frontendURL.JoinPath("/login/reset_password")
	q := resetURL.Query()
	q.Set("token", token)
	resetURL.RawQuery = q.Encode()

	msg := emailMessage.EmailMessage{
		EmailType: emailMessage.EmailTypeBuyerPasswordReset,
		To:        email,
		ResetURL:  resetURL.String(),
	}
	if err := u.publishEmail.Execute(ctx, msg); err != nil {
		return fmt.Errorf("failed to enqueue password reset email: %w", err)
	}

	return nil
}
