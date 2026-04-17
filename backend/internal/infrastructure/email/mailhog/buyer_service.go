package mailhog

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"

	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/email/templates"
)

var buyerSendMailFunc = smtp.SendMail

// BuyerEmailService provides BuyerEmailService related functionality.
type BuyerEmailService struct {
	cfg            config.EmailConfig
	templateLoader templates.TemplateProvider
}

var _ service.BuyerEmailService = (*BuyerEmailService)(nil)

// NewBuyerEmailService creates a new BuyerEmailService instance.
func NewBuyerEmailService(cfg config.EmailConfig, loader templates.TemplateProvider) *BuyerEmailService {
	return &BuyerEmailService{
		cfg:            cfg,
		templateLoader: loader,
	}
}

func (s *BuyerEmailService) send(to, subject, body string) error {
	msg := fmt.Appendf(nil, "To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s", to, subject, body)

	// MailHog doesn't require auth
	return buyerSendMailFunc(s.cfg.SMTPAddress(), nil, s.cfg.GetSMTPFrom(), []string{to}, msg)
}

// SendBuyerPasswordReset provides SendBuyerPasswordReset related functionality.
func (s *BuyerEmailService) SendBuyerPasswordReset(_ context.Context, to, url string) error {
	tmpl := s.templateLoader.Get("buyer_password_reset.txt")
	if tmpl == nil {
		return fmt.Errorf("template buyer_password_reset.txt not found")
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, map[string]string{"ResetURL": url}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	subject := "【Fish Auction】パスワード再設定のご案内"
	return s.send(to, subject, body.String())
}
