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

var adminSendMailFunc = smtp.SendMail

// AdminEmailService provides AdminEmailService related functionality.
type AdminEmailService struct {
	cfg            *config.Config
	templateLoader templates.TemplateProvider
}

var _ service.AdminEmailService = (*AdminEmailService)(nil)

// NewAdminEmailService creates a new AdminEmailService instance.
func NewAdminEmailService(cfg *config.Config, loader templates.TemplateProvider) *AdminEmailService {
	return &AdminEmailService{
		cfg:            cfg,
		templateLoader: loader,
	}
}

func (s *AdminEmailService) send(to, subject, body string) error {
	msg := fmt.Appendf(nil, "To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s", to, subject, body)

	// MailHog doesn't require auth
	return adminSendMailFunc(s.cfg.SMTPAddress(), nil, s.cfg.SMTPFrom, []string{to}, msg)
}

// SendAdminPasswordReset provides SendAdminPasswordReset related functionality.
func (s *AdminEmailService) SendAdminPasswordReset(_ context.Context, to, url string) error {
	tmpl := s.templateLoader.Get("admin_password_reset.txt")
	if tmpl == nil {
		return fmt.Errorf("template admin_password_reset.txt not found")
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, map[string]string{"ResetURL": url}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	subject := "【Fish Auction】パスワード再設定のご案内"
	return s.send(to, subject, body.String())
}
