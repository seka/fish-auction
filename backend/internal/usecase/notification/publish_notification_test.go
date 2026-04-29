package notification

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type mockPushService struct {
	publishFunc func(ctx context.Context, buyerID int, payload any) error
}

func (m *mockPushService) Send(_ context.Context, _ *model.PushSubscription, _ any) error {
	return nil
}

func (m *mockPushService) PublishToBuyer(ctx context.Context, buyerID int, payload any) error {
	return m.publishFunc(ctx, buyerID, payload)
}

func TestPublishNotificationUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	buyerID := 1
	payload := map[string]string{"title": "test", "body": "hello"}

	t.Run("success", func(t *testing.T) {
		publishCalled := false
		var capturedBuyerID int

		mockSvc := &mockPushService{
			publishFunc: func(_ context.Context, bID int, p any) error {
				publishCalled = true
				capturedBuyerID = bID
				return nil
			},
		}

		uc := NewPublishNotificationUseCase(mockSvc)
		err := uc.Execute(ctx, buyerID, payload)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !publishCalled {
			t.Error("Expected PushNotificationService.PublishToBuyer to be called")
		}
		if capturedBuyerID != buyerID {
			t.Errorf("Expected BuyerID %d, got %d", buyerID, capturedBuyerID)
		}
	})

	t.Run("publish error", func(t *testing.T) {
		publishErr := errors.New("publish failed")
		mockSvc := &mockPushService{
			publishFunc: func(_ context.Context, _ int, _ any) error {
				return publishErr
			},
		}

		uc := NewPublishNotificationUseCase(mockSvc)
		err := uc.Execute(ctx, buyerID, payload)

		if !errors.Is(err, publishErr) {
			t.Errorf("Expected error %v, got %v", publishErr, err)
		}
	})
}
