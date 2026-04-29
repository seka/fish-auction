package mailhog

import (
	"context"
	"errors"
	"net/smtp"
	"testing"

	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/infrastructure/email/templates"
)

func TestBuyerEmailService(t *testing.T) {
	// Setup real template loader for success cases
	realLoader, err := templates.NewTemplateLoader()
	if err != nil {
		t.Fatalf("failed to create template loader: %v", err)
	}

	cfg := &config.AppServerConfig{
		SMTPHost: "localhost",
		SMTPPort: "1025",
		SMTPFrom: "noreply@example.com",
	}

	t.Run("SendBuyerPasswordReset", func(t *testing.T) {
		tests := []struct {
			name        string
			to          string
			url         string
			mockSendErr error
			mockTmplErr bool
			wantErr     bool
		}{
			{
				name:    "Success",
				to:      "buyer@example.com",
				url:     "http://example.com/reset",
				wantErr: false,
			},
			{
				name:        "TemplateNotFound",
				to:          "buyer@example.com",
				url:         "http://example.com/reset",
				mockTmplErr: true,
				wantErr:     true,
			},
			{
				name:        "SendError",
				to:          "buyer@example.com",
				url:         "http://example.com/reset",
				mockSendErr: errors.New("smtp error"),
				wantErr:     true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				restore := setSendMailFunc(func(_ string, _ smtp.Auth, _ string, _ []string, _ []byte) error {
					if tt.mockSendErr != nil {
						return tt.mockSendErr
					}
					return nil
				})
				defer restore()

				loader := &mockTemplateLoader{realLoader: realLoader, mockErr: tt.mockTmplErr}
				svc := NewBuyerEmailService(cfg, loader)
				err := svc.SendBuyerPasswordReset(context.Background(), tt.to, tt.url)

				if (err != nil) != tt.wantErr {
					t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})
}
