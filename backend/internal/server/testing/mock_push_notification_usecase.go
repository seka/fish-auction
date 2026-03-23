package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockPushNotificationUseCase is a mock implementation of PushNotificationUseCase for testing.
type MockPushNotificationUseCase struct {
	SubscribeFunc        func(ctx context.Context, buyerID int, sub *model.PushSubscription) error
	SendNotificationFunc func(ctx context.Context, buyerID int, payload any) error
}

// Subscribe provides Subscribe related functionality.
func (m *MockPushNotificationUseCase) Subscribe(ctx context.Context, buyerID int, sub *model.PushSubscription) error {
	if m.SubscribeFunc != nil {
		return m.SubscribeFunc(ctx, buyerID, sub)
	}
	return nil
}

// SendNotification sends a message or notification.
func (m *MockPushNotificationUseCase) SendNotification(ctx context.Context, buyerID int, payload any) error {
	if m.SendNotificationFunc != nil {
		return m.SendNotificationFunc(ctx, buyerID, payload)
	}
	return nil
}
