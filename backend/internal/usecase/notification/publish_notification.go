package notification

import (
	"context"
	"errors"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

// PublishNotificationUseCase defines the interface for publish notifications.
type PublishNotificationUseCase interface {
	// Execute sends a notification to a buyer.
	Execute(ctx context.Context, buyerID int, payload any) error
}

type publishNotificationUseCase struct {
	repo                    repository.PushRepository
	pushNotificationService service.PushNotificationService
}

var _ PublishNotificationUseCase = (*publishNotificationUseCase)(nil)

// NewPublishNotificationUseCase creates a new instance of PublishNotificationUseCase.
func NewPublishNotificationUseCase(
	repo repository.PushRepository,
	pushNotificationService service.PushNotificationService,
) PublishNotificationUseCase {
	return &publishNotificationUseCase{
		repo:                    repo,
		pushNotificationService: pushNotificationService,
	}
}

func (uc *publishNotificationUseCase) Execute(ctx context.Context, buyerID int, payload any) error {
	subs, err := uc.repo.GetSubscriptionsByBuyerID(ctx, buyerID)
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		return nil
	}

	// Send to all subscriptions for the user
	// In production, this should probably be done asynchronously via a queue
	for _, sub := range subs {
		if err := uc.pushNotificationService.Send(ctx, &sub, payload); err != nil {
			// If subscription is expired or not found, delete it
			var goneErr *domainErrors.GoneError
			var notFoundErr *domainErrors.NotFoundError
			if (errors.As(err, &goneErr) && goneErr.Resource == "Subscription") ||
				(errors.As(err, &notFoundErr) && notFoundErr.Resource == "Subscription") {
				_ = uc.repo.DeleteSubscription(ctx, sub.Endpoint)
			}
			continue
		}
	}

	return nil
}
