package handler

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	notificationMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

type mockPushRepository struct {
	getSubscriptionsByBuyerIDFunc func(ctx context.Context, buyerID int) ([]model.PushSubscription, error)
	deleteSubscriptionFunc        func(ctx context.Context, endpoint string) error
}

func (m *mockPushRepository) GetSubscriptionsByBuyerID(ctx context.Context, buyerID int) ([]model.PushSubscription, error) {
	return m.getSubscriptionsByBuyerIDFunc(ctx, buyerID)
}

func (m *mockPushRepository) DeleteSubscription(ctx context.Context, endpoint string) error {
	return m.deleteSubscriptionFunc(ctx, endpoint)
}

func (m *mockPushRepository) SaveSubscription(_ context.Context, _ *model.PushSubscription) error {
	return nil
}

type mockPushNotificationService struct {
	sendFunc func(ctx context.Context, sub *model.PushSubscription, payload any) error
}

func (m *mockPushNotificationService) Send(ctx context.Context, sub *model.PushSubscription, payload any) error {
	return m.sendFunc(ctx, sub, payload)
}

func TestPushNotificationHandler_Handle_InvalidPayload(t *testing.T) {
	ctx := context.Background()

	h := NewPushNotificationHandler(&mockPushRepository{}, &mockPushNotificationService{})
	err := h.Handle(ctx, []byte("invalid json"))

	if err == nil {
		t.Error("Expected error for invalid JSON payload, got nil")
	}
}

func TestPushNotificationHandler_Handle_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("db connection failed")

	job := notificationMessage.PushNotificationMessage{BuyerID: 1, Payload: "test"}
	payloadBytes, _ := json.Marshal(job)

	repo := &mockPushRepository{
		getSubscriptionsByBuyerIDFunc: func(_ context.Context, _ int) ([]model.PushSubscription, error) {
			return nil, repoErr
		},
	}
	pushSvc := &mockPushNotificationService{
		sendFunc: func(_ context.Context, _ *model.PushSubscription, _ any) error {
			t.Error("Send should not be called")
			return nil
		},
	}

	h := NewPushNotificationHandler(repo, pushSvc)
	err := h.Handle(ctx, payloadBytes)

	if !errors.Is(err, repoErr) {
		t.Errorf("Expected error %v, got %v", repoErr, err)
	}
}

func TestPushNotificationHandler_Handle(t *testing.T) {
	ctx := context.Background()
	buyerID := 1
	jobPayload := map[string]string{"title": "test", "body": "hello"}
	job := notificationMessage.PushNotificationMessage{
		BuyerID: buyerID,
		Payload: jobPayload,
	}
	payloadBytes, _ := json.Marshal(job)

	t.Run("no subscriptions", func(t *testing.T) {
		repo := &mockPushRepository{
			getSubscriptionsByBuyerIDFunc: func(_ context.Context, _ int) ([]model.PushSubscription, error) {
				return []model.PushSubscription{}, nil
			},
		}
		pushSvc := &mockPushNotificationService{
			sendFunc: func(_ context.Context, _ *model.PushSubscription, _ any) error {
				t.Error("Send should not be called")
				return nil
			},
		}

		h := NewPushNotificationHandler(repo, pushSvc)
		err := h.Handle(ctx, payloadBytes)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("send success", func(t *testing.T) {
		sendCalled := 0
		repo := &mockPushRepository{
			getSubscriptionsByBuyerIDFunc: func(_ context.Context, _ int) ([]model.PushSubscription, error) {
				return []model.PushSubscription{
					{Endpoint: "endpoint1"},
					{Endpoint: "endpoint2"},
				}, nil
			},
		}
		pushSvc := &mockPushNotificationService{
			sendFunc: func(_ context.Context, _ *model.PushSubscription, _ any) error {
				sendCalled++
				return nil
			},
		}

		h := NewPushNotificationHandler(repo, pushSvc)
		err := h.Handle(ctx, payloadBytes)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if sendCalled != 2 {
			t.Errorf("Expected Send to be called 2 times, got %d", sendCalled)
		}
	})

	t.Run("cleanup on GoneError", func(t *testing.T) {
		deletedEndpoints := make([]string, 0)
		repo := &mockPushRepository{
			getSubscriptionsByBuyerIDFunc: func(_ context.Context, _ int) ([]model.PushSubscription, error) {
				return []model.PushSubscription{
					{Endpoint: "expired-endpoint"},
				}, nil
			},
			deleteSubscriptionFunc: func(_ context.Context, endpoint string) error {
				deletedEndpoints = append(deletedEndpoints, endpoint)
				return nil
			},
		}
		pushSvc := &mockPushNotificationService{
			sendFunc: func(_ context.Context, _ *model.PushSubscription, _ any) error {
				return &domainErrors.GoneError{Resource: "Subscription"}
			},
		}

		h := NewPushNotificationHandler(repo, pushSvc)
		err := h.Handle(ctx, payloadBytes)

		if err != nil {
			t.Errorf("Expected no error (handler should swallow individual push errors), got %v", err)
		}
		if len(deletedEndpoints) != 1 || deletedEndpoints[0] != "expired-endpoint" {
			t.Errorf("Expected DeleteSubscription to be called for expired-endpoint, got %v", deletedEndpoints)
		}
	})

	t.Run("cleanup on NotFoundError", func(t *testing.T) {
		deletedEndpoints := make([]string, 0)
		repo := &mockPushRepository{
			getSubscriptionsByBuyerIDFunc: func(_ context.Context, _ int) ([]model.PushSubscription, error) {
				return []model.PushSubscription{
					{Endpoint: "notfound-endpoint"},
				}, nil
			},
			deleteSubscriptionFunc: func(_ context.Context, endpoint string) error {
				deletedEndpoints = append(deletedEndpoints, endpoint)
				return nil
			},
		}
		pushSvc := &mockPushNotificationService{
			sendFunc: func(_ context.Context, _ *model.PushSubscription, _ any) error {
				return &domainErrors.NotFoundError{Resource: "Subscription"}
			},
		}

		h := NewPushNotificationHandler(repo, pushSvc)
		err := h.Handle(ctx, payloadBytes)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(deletedEndpoints) != 1 || deletedEndpoints[0] != "notfound-endpoint" {
			t.Errorf("Expected DeleteSubscription to be called for notfound-endpoint, got %v", deletedEndpoints)
		}
	})

	t.Run("no cleanup on other errors", func(t *testing.T) {
		deletedEndpoints := make([]string, 0)
		repo := &mockPushRepository{
			getSubscriptionsByBuyerIDFunc: func(_ context.Context, _ int) ([]model.PushSubscription, error) {
				return []model.PushSubscription{
					{Endpoint: "temp-error-endpoint"},
				}, nil
			},
			deleteSubscriptionFunc: func(_ context.Context, endpoint string) error {
				deletedEndpoints = append(deletedEndpoints, endpoint)
				return nil
			},
		}
		pushSvc := &mockPushNotificationService{
			sendFunc: func(_ context.Context, _ *model.PushSubscription, _ any) error {
				return errors.New("temporary error")
			},
		}

		h := NewPushNotificationHandler(repo, pushSvc)
		err := h.Handle(ctx, payloadBytes)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(deletedEndpoints) != 0 {
			t.Errorf("Expected DeleteSubscription NOT to be called, but got %v", deletedEndpoints)
		}
	})
}
