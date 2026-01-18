package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type MockPushNotificationUseCase struct {
	SubscribeFunc        func(ctx context.Context, buyerID int, sub *model.PushSubscription) error
	SendNotificationFunc func(ctx context.Context, buyerID int, payload interface{}) error
}

func (m *MockPushNotificationUseCase) Subscribe(ctx context.Context, buyerID int, sub *model.PushSubscription) error {
	if m.SubscribeFunc != nil {
		return m.SubscribeFunc(ctx, buyerID, sub)
	}
	return nil
}

func (m *MockPushNotificationUseCase) SendNotification(ctx context.Context, buyerID int, payload interface{}) error {
	if m.SendNotificationFunc != nil {
		return m.SendNotificationFunc(ctx, buyerID, payload)
	}
	return nil
}
