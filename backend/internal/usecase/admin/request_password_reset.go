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
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	emailMessage "github.com/seka/fish-auction/backend/internal/event"
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
	jobQueue     service.JobQueue
	frontendURL  *url.URL
	txMgr        repository.TransactionManager
	clock        service.Clock
}

var _ RequestPasswordResetUseCase = (*requestPasswordResetUseCase)(nil)

// NewRequestPasswordResetUseCase creates a new RequestPasswordResetUseCase instance.
func NewRequestPasswordResetUseCase(
	adminRepo repository.AdminRepository,
	pwdResetRepo repository.PasswordResetRepository,
	jobQueue service.JobQueue,
	frontendURL *url.URL,
	txMgr repository.TransactionManager,
	clock service.Clock,
) RequestPasswordResetUseCase {
	return &requestPasswordResetUseCase{
		adminRepo:    adminRepo,
		pwdResetRepo: pwdResetRepo,
		jobQueue:     jobQueue,
		frontendURL:  frontendURL,
		txMgr:        txMgr,
		clock:        clock,
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

	// 3. Save token in DB transaction (expires in 30 mins)
	resetURL := u.frontendURL.JoinPath("/login/admin/reset_password")
	q := resetURL.Query()
	q.Set("token", token)
	resetURL.RawQuery = q.Encode()

	expiresAt := u.clock.Now().Add(30 * time.Minute)
	if err := u.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := u.pwdResetRepo.DeleteAllByUserID(txCtx, admin.ID, "admin"); err != nil {
			return fmt.Errorf("failed to invalidate old reset tokens: %w", err)
		}
		if err = u.pwdResetRepo.Create(txCtx, admin.ID, "admin", tokenHash, expiresAt); err != nil {
			return fmt.Errorf("failed to create new reset token: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}

	wire := emailMessage.EmailMessage{
		EmailType: emailMessage.EmailTypeAdminPasswordReset,
		To:        email,
		ResetURL:  resetURL.String(),
	}
	if err := u.jobQueue.Enqueue(ctx, model.JobTypeEmail, wire); err != nil {
		return fmt.Errorf("failed to enqueue password reset email: %w", err)
	}

	return nil
}
