package admin

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

type RequestPasswordResetUseCase interface {
	Execute(ctx context.Context, email string) error
}

type requestPasswordResetUseCase struct {
	adminRepo    repository.AdminRepository
	pwdResetRepo repository.AdminPasswordResetRepository
	emailService service.EmailService
}

func NewRequestPasswordResetUseCase(
	adminRepo repository.AdminRepository,
	pwdResetRepo repository.AdminPasswordResetRepository,
	emailService service.EmailService,
) RequestPasswordResetUseCase {
	return &requestPasswordResetUseCase{
		adminRepo:    adminRepo,
		pwdResetRepo: pwdResetRepo,
		emailService: emailService,
	}
}

func (u *requestPasswordResetUseCase) Execute(ctx context.Context, email string) error {
	admin, err := u.adminRepo.FindOneByEmail(ctx, email)
	if err != nil {
		return nil
	}
	if admin == nil {
		return nil
	}

	// 1. Generate secure token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// 2. Hash token for DB
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// 3. Save to DB (expires in 30 mins)
	expiresAt := time.Now().Add(30 * time.Minute)
	if err := u.pwdResetRepo.DeleteAllByAdminID(ctx, admin.ID); err != nil {
		return fmt.Errorf("failed to delete old tokens: %w", err)
	}
	if err := u.pwdResetRepo.Create(ctx, admin.ID, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// 4. Send Email
	resetURL := fmt.Sprintf("http://localhost:3000/login/admin/reset_password?token=%s", token) // TODO: config
	if err := u.emailService.SendAdminPasswordReset(ctx, email, resetURL); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
