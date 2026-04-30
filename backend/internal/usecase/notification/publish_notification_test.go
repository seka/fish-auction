package notification

import (
	"context"
	"errors"
	"testing"
)

type mockPushQueue struct {
	enqueueFunc func(ctx context.Context, buyerID int, payload any) error
}

func (m *mockPushQueue) Enqueue(ctx context.Context, buyerID int, payload any) error {
	return m.enqueueFunc(ctx, buyerID, payload)
}

func TestPublishNotificationUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	buyerID := 1
	payload := map[string]string{"title": "test", "body": "hello"}

	t.Run("success", func(t *testing.T) {
		enqueueCalled := false
		var capturedBuyerID int

		mockQueue := &mockPushQueue{
			enqueueFunc: func(_ context.Context, bID int, _ any) error {
				enqueueCalled = true
				capturedBuyerID = bID
				return nil
			},
		}

		uc := NewPublishNotificationUseCase(mockQueue)
		err := uc.Execute(ctx, buyerID, payload)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !enqueueCalled {
			t.Error("Expected PushNotificationQueue.Enqueue to be called")
		}
		if capturedBuyerID != buyerID {
			t.Errorf("Expected BuyerID %d, got %d", buyerID, capturedBuyerID)
		}
	})

	t.Run("enqueue error", func(t *testing.T) {
		enqueueErr := errors.New("enqueue failed")
		mockQueue := &mockPushQueue{
			enqueueFunc: func(_ context.Context, _ int, _ any) error {
				return enqueueErr
			},
		}

		uc := NewPublishNotificationUseCase(mockQueue)
		err := uc.Execute(ctx, buyerID, payload)

		if !errors.Is(err, enqueueErr) {
			t.Errorf("Expected error %v, got %v", enqueueErr, err)
		}
	})
}
