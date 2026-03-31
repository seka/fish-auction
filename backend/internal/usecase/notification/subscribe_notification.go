package notification

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// SubscribeNotificationUseCase defines the interface for subscribe notifications.
type SubscribeNotificationUseCase interface {
	// Execute saves a push subscription for a buyer.
	Execute(ctx context.Context, buyerID int, sub *model.PushSubscription) error
}

type subscribeNotificationUseCase struct {
	repo repository.PushRepository
}

var _ SubscribeNotificationUseCase = (*subscribeNotificationUseCase)(nil)

// NewSubscribeNotificationUseCase creates a new instance of SubscribeNotificationUseCase.
func NewSubscribeNotificationUseCase(repo repository.PushRepository) SubscribeNotificationUseCase {
	return &subscribeNotificationUseCase{
		repo: repo,
	}
}

func (uc *subscribeNotificationUseCase) Execute(ctx context.Context, buyerID int, sub *model.PushSubscription) error {
	sub.BuyerID = buyerID
	return uc.repo.SaveSubscription(ctx, sub)
}
