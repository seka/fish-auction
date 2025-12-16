package admin

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"time"

	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type RequestPasswordResetUseCase interface {
	Execute(ctx context.Context, email string) error
}

type requestPasswordResetUseCase struct {
	adminRepo    repository.AdminRepository
	pwdResetRepo repository.AdminPasswordResetRepository
	cfg          *config.Config
}

func NewRequestPasswordResetUseCase(
	adminRepo repository.AdminRepository,
	pwdResetRepo repository.AdminPasswordResetRepository,
	cfg *config.Config,
) RequestPasswordResetUseCase {
	return &requestPasswordResetUseCase{
		adminRepo:    adminRepo,
		pwdResetRepo: pwdResetRepo,
		cfg:          cfg,
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
	if err := u.sendEmail(email, token); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (u *requestPasswordResetUseCase) sendEmail(to, token string) error {
	resetURL := fmt.Sprintf("%s/login/admin/reset_password?token=%s", "http://localhost:3000", token) // TODO: config

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Admin Password Reset Request\r\n"+
		"\r\n"+
		"Click the link below to reset your password:\r\n"+
		"%s\r\n"+
		"\r\n"+
		"This link expires in 30 minutes.\r\n", to, resetURL))

	addr := fmt.Sprintf("%s:%s", u.cfg.SMTPHost, u.cfg.SMTPPort)
	return smtp.SendMail(addr, nil, u.cfg.SMTPFrom, []string{to}, msg)
}
