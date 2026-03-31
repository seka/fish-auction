package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockSubscribeNotificationUseCase is a mock implementation of SubscribeNotificationUseCase for testing.
type MockSubscribeNotificationUseCase struct {
	ExecuteFunc func(ctx context.Context, buyerID int, sub *model.PushSubscription) error
}

func (m *MockSubscribeNotificationUseCase) Execute(ctx context.Context, buyerID int, sub *model.PushSubscription) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, buyerID, sub)
	}
	return nil
}

// MockPublishNotificationUseCase is a mock implementation of PublishNotificationUseCase for testing.
type MockPublishNotificationUseCase struct {
	ExecuteFunc func(ctx context.Context, buyerID int, payload any) error
}

func (m *MockPublishNotificationUseCase) Execute(ctx context.Context, buyerID int, payload any) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, buyerID, payload)
	}
	return nil
}
