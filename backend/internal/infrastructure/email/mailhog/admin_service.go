package mailhog

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"

	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/infrastructure/email/templates"
)

var adminSendMailFunc = smtp.SendMail

type AdminEmailService struct {
	cfg            *config.Config
	templateLoader templates.TemplateProvider
}

func NewAdminEmailService(cfg *config.Config, loader templates.TemplateProvider) *AdminEmailService {
	return &AdminEmailService{
		cfg:            cfg,
		templateLoader: loader,
	}
}

func (s *AdminEmailService) send(to, subject, body string) error {
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s", to, subject, body))

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	// MailHog doesn't require auth
	return adminSendMailFunc(addr, nil, s.cfg.SMTPFrom, []string{to}, msg)
}

func (s *AdminEmailService) SendAdminPasswordReset(ctx context.Context, to, url string) error {
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
