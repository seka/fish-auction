package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	emailMessage "github.com/seka/fish-auction/backend/internal/event"
)

type emailHandler struct {
	buyerEmailSvc service.BuyerEmailService
	adminEmailSvc service.AdminEmailService
}

// NewEmailHandler creates a new handler for email jobs.
func NewEmailHandler(buyerEmailSvc service.BuyerEmailService, adminEmailSvc service.AdminEmailService) *emailHandler {
	return &emailHandler{
		buyerEmailSvc: buyerEmailSvc,
		adminEmailSvc: adminEmailSvc,
	}
}

func (h *emailHandler) Handle(ctx context.Context, msg *model.JobMessage) error {
	var emailMsg emailMessage.EmailMessage
	if err := json.Unmarshal(msg.Payload, &emailMsg); err != nil {
		return fmt.Errorf("failed to unmarshal email job payload: %w", err)
	}

	switch emailMsg.EmailType {
	case emailMessage.EmailTypeBuyerPasswordReset:
		return h.buyerEmailSvc.SendBuyerPasswordReset(ctx, emailMsg.To, emailMsg.ResetURL)
	case emailMessage.EmailTypeAdminPasswordReset:
		return h.adminEmailSvc.SendAdminPasswordReset(ctx, emailMsg.To, emailMsg.ResetURL)
	default:
		return fmt.Errorf("unsupported email type: %s", emailMsg.EmailType)
	}
}
