package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	notificationMessage "github.com/seka/fish-auction/backend/internal/event"
)

// PublishNotificationUseCase defines the interface for publish notifications.
type PublishNotificationUseCase interface {
	// Execute enqueues a notification job for a buyer.
	Execute(ctx context.Context, buyerID int, payload any) error
}

type publishNotificationUseCase struct {
	outboxRepo repository.OutboxRepository
}

var _ PublishNotificationUseCase = (*publishNotificationUseCase)(nil)

// NewPublishNotificationUseCase creates a new instance of PublishNotificationUseCase.
func NewPublishNotificationUseCase(
	outboxRepo repository.OutboxRepository,
) PublishNotificationUseCase {
	return &publishNotificationUseCase{
		outboxRepo: outboxRepo,
	}
}

func (uc *publishNotificationUseCase) Execute(ctx context.Context, buyerID int, payload any) error {
	wire := notificationMessage.PushNotificationMessage{
		BuyerID: buyerID,
		Payload: payload,
	}
	body, err := json.Marshal(wire)
	if err != nil {
		return fmt.Errorf("failed to marshal push notification payload: %w", err)
	}
	return uc.outboxRepo.Insert(ctx, model.JobTypePushNotification, 1, body)
}
