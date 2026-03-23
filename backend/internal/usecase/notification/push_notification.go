package notification

import (
	"context"
	"strings"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

// PushNotificationUseCase defines the interface for push notifications.
type PushNotificationUseCase interface {
	// Subscribe saves a push subscription for a buyer.
	Subscribe(ctx context.Context, buyerID int, sub *model.PushSubscription) error
	// SendNotification sends a notification to a buyer.
	SendNotification(ctx context.Context, buyerID int, payload any) error
}

type pushNotificationUseCase struct {
	repo                    repository.PushRepository
	pushNotificationService service.PushNotificationService
}

var _ PushNotificationUseCase = (*pushNotificationUseCase)(nil)

// NewPushNotificationUseCase creates a new instance of PushNotificationUseCase.
func NewPushNotificationUseCase(
	repo repository.PushRepository,
	pushNotificationService service.PushNotificationService,
) PushNotificationUseCase {
	return &pushNotificationUseCase{
		repo:                    repo,
		pushNotificationService: pushNotificationService,
	}
}

func (uc *pushNotificationUseCase) Subscribe(ctx context.Context, buyerID int, sub *model.PushSubscription) error {
	sub.BuyerID = buyerID
	return uc.repo.SaveSubscription(ctx, sub)
}

func (uc *pushNotificationUseCase) SendNotification(ctx context.Context, buyerID int, payload any) error {
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
			if strings.Contains(err.Error(), "expired") || strings.Contains(err.Error(), "status 410") || strings.Contains(err.Error(), "status 404") {
				_ = uc.repo.DeleteSubscription(ctx, sub.Endpoint)
			}
			continue
		}
	}

	return nil
}
