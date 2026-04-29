package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/job/handler"
)

type mockBuyerEmailSvc struct {
	err error
}

func (m *mockBuyerEmailSvc) SendBuyerPasswordReset(_ context.Context, _, _ string) error {
	return m.err
}

type mockAdminEmailSvc struct {
	err error
}

func (m *mockAdminEmailSvc) SendAdminPasswordReset(_ context.Context, _, _ string) error {
	return m.err
}

func TestEmailHandler_Handle(t *testing.T) {
	buyerPayload := `{"email_type":"buyer_password_reset","to":"buyer@example.com","reset_url":"https://example.com/reset"}`
	adminPayload := `{"email_type":"admin_password_reset","to":"admin@example.com","reset_url":"https://example.com/admin/reset"}`

	tests := []struct {
		name     string
		payload  string
		buyerErr error
		adminErr error
		wantErr  bool
	}{
		{
			name:    "buyer password reset success",
			payload: buyerPayload,
		},
		{
			name:    "admin password reset success",
			payload: adminPayload,
		},
		{
			name:     "buyer email service error",
			payload:  buyerPayload,
			buyerErr: errors.New("smtp error"),
			wantErr:  true,
		},
		{
			name:     "admin email service error",
			payload:  adminPayload,
			adminErr: errors.New("smtp error"),
			wantErr:  true,
		},
		{
			name:    "invalid payload",
			payload: `{invalid json`,
			wantErr: true,
		},
		{
			name:    "unknown email type",
			payload: `{"email_type":"unknown_type","to":"x@example.com"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewEmailHandler(
				&mockBuyerEmailSvc{err: tt.buyerErr},
				&mockAdminEmailSvc{err: tt.adminErr},
			)
			err := h.Handle(context.Background(), []byte(tt.payload))
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}
