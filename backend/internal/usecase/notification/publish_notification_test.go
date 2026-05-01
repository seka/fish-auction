package notification

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type mockPushQueue struct {
	enqueueFunc func(ctx context.Context, jobType model.JobType, payload any) error
}

func (m *mockPushQueue) Enqueue(ctx context.Context, jobType model.JobType, payload any) error {
	return m.enqueueFunc(ctx, jobType, payload)
}

func (m *mockPushQueue) EnqueueRaw(_ context.Context, _ model.JobType, _ []byte) error {
	return nil
}

func (m *mockPushQueue) Dequeue(_ context.Context, _ int32) ([]*model.JobMessage, error) {
	return nil, nil
}

func (m *mockPushQueue) DeleteMessage(_ context.Context, _ *model.JobMessage) error {
	return nil
}

func TestPublishNotificationUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	buyerID := 1
	payload := map[string]string{"title": "test", "body": "hello"}

	t.Run("success", func(t *testing.T) {
		enqueueCalled := false
		var capturedJobType model.JobType

		mockQueue := &mockPushQueue{
			enqueueFunc: func(_ context.Context, jt model.JobType, _ any) error {
				enqueueCalled = true
				capturedJobType = jt
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
		if capturedJobType != model.JobTypePushNotification {
			t.Errorf("Expected JobType %s, got %s", model.JobTypePushNotification, capturedJobType)
		}
	})

	t.Run("enqueue error", func(t *testing.T) {
		enqueueErr := errors.New("enqueue failed")
		mockQueue := &mockPushQueue{
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
