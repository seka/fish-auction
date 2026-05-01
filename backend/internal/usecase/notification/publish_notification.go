package notification

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
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
	return uc.outboxRepo.InsertPushNotificationJob(ctx, buyerID, payload)
}
