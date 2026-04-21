package notification

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	notificationMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

type mockJobQueue struct {
	enqueueFunc func(ctx context.Context, jobType model.JobType, payload any) error
}

func (m *mockJobQueue) Enqueue(ctx context.Context, jobType model.JobType, payload any) error {
	return m.enqueueFunc(ctx, jobType, payload)
}

func (m *mockJobQueue) Dequeue(_ context.Context, _ int32) ([]*model.JobMessage, error) {
	return nil, nil
}

func (m *mockJobQueue) DeleteMessage(_ context.Context, _ *model.JobMessage) error {
	return nil
}

func TestPublishNotificationUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	buyerID := 1
	payload := map[string]string{"title": "test", "body": "hello"}

	t.Run("success", func(t *testing.T) {
		enqueueCalled := false
		var capturedJobType model.JobType
		var capturedPayload any

		mockQueue := &mockJobQueue{
			enqueueFunc: func(_ context.Context, jobType model.JobType, payload any) error {
				enqueueCalled = true
				capturedJobType = jobType
				capturedPayload = payload
				return nil
			},
		}

		uc := NewPublishNotificationUseCase(mockQueue)
		err := uc.Execute(ctx, buyerID, payload)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !enqueueCalled {
			t.Error("Expected JobQueue.Enqueue to be called")
		}
		if capturedJobType != model.JobTypePushNotification {
			t.Errorf("Expected jobType '%s', got '%s'", model.JobTypePushNotification, capturedJobType)
		}

		p, ok := capturedPayload.(notificationMessage.PushNotificationMessage)
		if !ok {
			t.Fatalf("Captured payload is not the expected struct type: %T", capturedPayload)
		}
		if p.BuyerID != buyerID {
			t.Errorf("Expected BuyerID %d, got %d", buyerID, p.BuyerID)
		}
	})

	t.Run("enqueue error", func(t *testing.T) {
		enqueueErr := errors.New("enqueue failed")
		mockQueue := &mockJobQueue{
			enqueueFunc: func(_ context.Context, _ model.JobType, _ any) error {
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
