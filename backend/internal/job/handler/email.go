package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/seka/fish-auction/backend/internal/domain/service"
	emailMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

type emailHandler struct {
	buyerEmailSvc service.BuyerEmailService
	adminEmailSvc service.AdminEmailService
}

// NewEmailHandler creates a new handler for email jobs.
func NewEmailHandler(buyerEmailSvc service.BuyerEmailService, adminEmailSvc service.AdminEmailService) Handler {
	return &emailHandler{
		buyerEmailSvc: buyerEmailSvc,
		adminEmailSvc: adminEmailSvc,
	}
}

func (h *emailHandler) Handle(ctx context.Context, payload []byte) error {
	var msg emailMessage.EmailMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal email job payload: %w", err)
	}

	switch msg.EmailType {
	case emailMessage.EmailTypeBuyerPasswordReset:
		return h.buyerEmailSvc.SendBuyerPasswordReset(ctx, msg.To, msg.Data["ResetURL"])
	case emailMessage.EmailTypeAdminPasswordReset:
		return h.adminEmailSvc.SendAdminPasswordReset(ctx, msg.To, msg.Data["ResetURL"])
	default:
		return fmt.Errorf("unsupported email type: %s", msg.EmailType)
	}
}
