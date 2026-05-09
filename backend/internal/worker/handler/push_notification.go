package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	notificationMessage "github.com/seka/fish-auction/backend/internal/event"
)

// pushNotificationHandler implements the Handler interface for push notifications.
type pushNotificationHandler struct {
	repo    repository.PushRepository
	pushSvc service.PushNotificationService
}

// NewPushNotificationHandler creates a new handler for push notification jobs.
func NewPushNotificationHandler(
	repo repository.PushRepository,
	pushSvc service.PushNotificationService,
) *pushNotificationHandler {
	return &pushNotificationHandler{
		repo:    repo,
		pushSvc: pushSvc,
	}
}

func (h *pushNotificationHandler) Handle(ctx context.Context, msg *model.JobMessage) error {
	var job notificationMessage.PushNotificationMessage
	if err := json.Unmarshal(msg.Payload, &job); err != nil {
		return fmt.Errorf("failed to unmarshal job payload: %w", err)
	}

	// 1. Get subscriptions
	subs, err := h.repo.GetSubscriptionsByBuyerID(ctx, job.BuyerID)
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	if len(subs) == 0 {
		return nil
	}

	// 2. Send notifications
	for _, sub := range subs {
		if err := h.pushSvc.Send(ctx, &sub, job.Payload); err != nil {
			slog.Error("failed to send push notification", "buyer_id", job.BuyerID, "endpoint", sub.Endpoint, "err", err)

			// Cleanup expired subscriptions
			var goneErr *domainErrors.GoneError
			var notFoundErr *domainErrors.NotFoundError
			if (errors.As(err, &goneErr) && goneErr.Resource == "Subscription") ||
				(errors.As(err, &notFoundErr) && notFoundErr.Resource == "Subscription") {
				_ = h.repo.DeleteSubscription(ctx, sub.Endpoint)
			}
		}
	}

	return nil
}
