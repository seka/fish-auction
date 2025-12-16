package auth

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
	buyerRepo    repository.BuyerRepository
	pwdResetRepo repository.BuyerPasswordResetRepository
	cfg          *config.Config
}

func NewRequestPasswordResetUseCase(
	buyerRepo repository.BuyerRepository,
	pwdResetRepo repository.BuyerPasswordResetRepository,
	cfg *config.Config,
) RequestPasswordResetUseCase {
	return &requestPasswordResetUseCase{
		buyerRepo:    buyerRepo,
		pwdResetRepo: pwdResetRepo,
		cfg:          cfg,
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
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// 2. Hash token for DB
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// 3. Save to DB (expires in 30 mins)
	expiresAt := time.Now().Add(30 * time.Minute)
	// Invalidate old tokens for this user first
	if err := u.pwdResetRepo.DeleteAllByBuyerID(ctx, buyer.ID); err != nil {
		return fmt.Errorf("failed to delete old tokens: %w", err)
	}
	if err := u.pwdResetRepo.Create(ctx, buyer.ID, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// 4. Send Email
	if err := u.sendEmail(email, token); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (u *requestPasswordResetUseCase) sendEmail(to, token string) error {
	resetURL := fmt.Sprintf("%s/login/reset_password?token=%s", "http://localhost:3000", token) // TODO: frontend base url from config

	subject := "【Fish Auction】パスワード再設定のご案内"
	body := fmt.Sprintf(`
いつもFish Auctionをご利用いただきありがとうございます。
パスワード再設定のリクエストを受け付けました。

以下のリンクをクリックして、新しいパスワードを設定してください。

%s

※このリンクは30分間有効です。
※本メールに心当たりがない場合は、破棄してください。

--------------------------------------------------
Fish Auction 運営事務局
--------------------------------------------------
`, resetURL)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s", to, subject, body))

	addr := fmt.Sprintf("%s:%s", u.cfg.SMTPHost, u.cfg.SMTPPort)
	return smtp.SendMail(addr, nil, u.cfg.SMTPFrom, []string{to}, msg)
}
