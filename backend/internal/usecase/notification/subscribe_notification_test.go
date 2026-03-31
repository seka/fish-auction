package notification

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type mockPushRepository struct {
	saveSubscriptionFunc          func(ctx context.Context, sub *model.PushSubscription) error
	getSubscriptionsByBuyerIDFunc func(ctx context.Context, buyerID int) ([]model.PushSubscription, error)
	deleteSubscriptionFunc        func(ctx context.Context, endpoint string) error
}

func (m *mockPushRepository) SaveSubscription(ctx context.Context, sub *model.PushSubscription) error {
	return m.saveSubscriptionFunc(ctx, sub)
}

func (m *mockPushRepository) GetSubscriptionsByBuyerID(ctx context.Context, buyerID int) ([]model.PushSubscription, error) {
	return m.getSubscriptionsByBuyerIDFunc(ctx, buyerID)
}

func (m *mockPushRepository) DeleteSubscription(ctx context.Context, endpoint string) error {
	return m.deleteSubscriptionFunc(ctx, endpoint)
}

func TestSubscribeNotificationUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	buyerID := 1
	sub := &model.PushSubscription{
		Endpoint: "test-endpoint",
	}

	t.Run("success", func(t *testing.T) {
		repo := &mockPushRepository{
			saveSubscriptionFunc: func(_ context.Context, s *model.PushSubscription) error {
				if s.BuyerID != buyerID {
					t.Errorf("expected buyerID %d, got %d", buyerID, s.BuyerID)
				}
				if s.Endpoint != sub.Endpoint {
					t.Errorf("expected endpoint %s, got %s", sub.Endpoint, s.Endpoint)
				}
				return nil
			},
		}

		uc := NewSubscribeNotificationUseCase(repo)
		err := uc.Execute(ctx, buyerID, sub)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("error from repository", func(t *testing.T) {
		expectedErr := errors.New("repository error")
		repo := &mockPushRepository{
			saveSubscriptionFunc: func(_ context.Context, _ *model.PushSubscription) error {
				return expectedErr
			},
		}

		uc := NewSubscribeNotificationUseCase(repo)
		err := uc.Execute(ctx, buyerID, sub)
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}
