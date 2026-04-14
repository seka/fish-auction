package notification

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type mockJobQueue struct {
	enqueueFunc func(ctx context.Context, jobType model.JobType, payload []byte) error
}

func (m *mockJobQueue) Enqueue(ctx context.Context, jobType model.JobType, payload []byte) error {
	return m.enqueueFunc(ctx, jobType, payload)
}

func TestPublishNotificationUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	buyerID := 1
	payload := map[string]string{"title": "test", "body": "hello"}

	t.Run("success", func(t *testing.T) {
		enqueueCalled := false
		var capturedJobType model.JobType
		var capturedPayload []byte

		mockQueue := &mockJobQueue{
			enqueueFunc: func(_ context.Context, jobType model.JobType, payload []byte) error {
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

		type pushNotificationJobDTO struct {
			BuyerID int `json:"buyer_id"`
			Payload any `json:"payload"`
		}

		var job pushNotificationJobDTO
		if err := json.Unmarshal(capturedPayload, &job); err != nil {
			t.Fatalf("Failed to unmarshal captured payload: %v", err)
		}
		if job.BuyerID != buyerID {
			t.Errorf("Expected BuyerID %d, got %d", buyerID, job.BuyerID)
		}
	})

	t.Run("enqueue error", func(t *testing.T) {
		enqueueErr := errors.New("enqueue failed")
		mockQueue := &mockJobQueue{
			enqueueFunc: func(_ context.Context, _ model.JobType, _ []byte) error {
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
