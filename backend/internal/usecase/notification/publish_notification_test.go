package notification

import (
	"context"
	"errors"
	"testing"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type mockPushNotificationService struct {
	sendFunc func(ctx context.Context, sub *model.PushSubscription, payload any) error
}

func (m *mockPushNotificationService) Send(ctx context.Context, sub *model.PushSubscription, payload any) error {
	return m.sendFunc(ctx, sub, payload)
}

func TestPublishNotificationUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	buyerID := 1
	payload := "test-notification"

	t.Run("success", func(t *testing.T) {
		subs := []model.PushSubscription{
			{Endpoint: "endpoint-1"},
			{Endpoint: "endpoint-2"},
		}

		repo := &mockPushRepository{
			getSubscriptionsByBuyerIDFunc: func(_ context.Context, _ int) ([]model.PushSubscription, error) {
				return subs, nil
			},
		}

		sendCount := 0
		svc := &mockPushNotificationService{
			sendFunc: func(_ context.Context, _ *model.PushSubscription, _ any) error {
				sendCount++
				return nil
			},
		}

		uc := NewPublishNotificationUseCase(repo, svc)
		err := uc.Execute(ctx, buyerID, payload)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if sendCount != len(subs) {
			t.Errorf("Expected %d sends, got %d", len(subs), sendCount)
		}
	})

	t.Run("no subscriptions", func(t *testing.T) {
		repo := &mockPushRepository{
			getSubscriptionsByBuyerIDFunc: func(_ context.Context, _ int) ([]model.PushSubscription, error) {
				return []model.PushSubscription{}, nil
			},
		}

		svc := &mockPushNotificationService{
			sendFunc: func(_ context.Context, _ *model.PushSubscription, _ any) error {
				t.Error("Send should not be called")
				return nil
			},
		}

		uc := NewPublishNotificationUseCase(repo, svc)
		err := uc.Execute(ctx, buyerID, payload)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("error from repository", func(t *testing.T) {
		repoErr := errors.New("repository error")
		repo := &mockPushRepository{
			getSubscriptionsByBuyerIDFunc: func(_ context.Context, _ int) ([]model.PushSubscription, error) {
				return nil, repoErr
			},
		}

		uc := NewPublishNotificationUseCase(repo, &mockPushNotificationService{})
		err := uc.Execute(ctx, buyerID, payload)
		if !errors.Is(err, repoErr) {
			t.Errorf("Expected error %v, got %v", repoErr, err)
		}
	})

	t.Run("cleanup invalid subscriptions", func(t *testing.T) {
		subs := []model.PushSubscription{
			{Endpoint: "endpoint-1"},
			{Endpoint: "endpoint-2"},
		}

		deletedEndpoints := make(map[string]bool)
		repo := &mockPushRepository{
			getSubscriptionsByBuyerIDFunc: func(_ context.Context, _ int) ([]model.PushSubscription, error) {
				return subs, nil
			},
			deleteSubscriptionFunc: func(_ context.Context, endpoint string) error {
				deletedEndpoints[endpoint] = true
				return nil
			},
		}

		svc := &mockPushNotificationService{
			sendFunc: func(_ context.Context, sub *model.PushSubscription, _ any) error {
				if sub.Endpoint == "endpoint-1" {
					return &domainErrors.GoneError{Resource: "Subscription"}
				}
				if sub.Endpoint == "endpoint-2" {
					return &domainErrors.NotFoundError{Resource: "Subscription"}
				}
				return nil
			},
		}

		uc := NewPublishNotificationUseCase(repo, svc)
		err := uc.Execute(ctx, buyerID, payload)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !deletedEndpoints["endpoint-1"] {
			t.Error("Expected endpoint-1 to be deleted")
		}
		if !deletedEndpoints["endpoint-2"] {
			t.Error("Expected endpoint-2 to be deleted")
		}
	})
}
