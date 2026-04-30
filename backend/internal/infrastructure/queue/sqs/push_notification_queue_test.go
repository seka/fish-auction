package sqs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/mock"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/sqs"
	notificationMessage "github.com/seka/fish-auction/backend/internal/worker/message"
)

func TestPushNotificationQueue_Enqueue(t *testing.T) {
	ctx := context.Background()
	buyerID := 1
	payload := map[string]string{"title": "test", "body": "hello"}

	t.Run("success", func(t *testing.T) {
		enqueueCalled := false
		var capturedJobType model.JobType
		var capturedPayload any

		mockQueue := &mock.MockJobQueue{
			EnqueueFunc: func(_ context.Context, jobType model.JobType, payload any) error {
				enqueueCalled = true
				capturedJobType = jobType
				capturedPayload = payload
				return nil
			},
		}

		svc := sqs.NewPushNotificationQueue(mockQueue)
		err := svc.Enqueue(ctx, buyerID, payload)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !enqueueCalled {
			t.Fatal("Expected JobQueue.Enqueue to be called")
		}
		if capturedJobType != model.JobTypePushNotification {
			t.Errorf("Expected jobType '%s', got '%s'", model.JobTypePushNotification, capturedJobType)
		}

		p, ok := capturedPayload.(notificationMessage.PushNotificationMessage)
		if !ok {
			t.Fatalf("Captured payload is not of type PushNotificationMessage: %T", capturedPayload)
		}
		if p.BuyerID != buyerID {
			t.Errorf("Expected BuyerID %d, got %d", buyerID, p.BuyerID)
		}
	})

	t.Run("enqueue error", func(t *testing.T) {
		expectedErr := errors.New("enqueue error")
		mockQueue := &mock.MockJobQueue{
			EnqueueFunc: func(_ context.Context, _ model.JobType, _ any) error {
				return expectedErr
			},
		}

		svc := sqs.NewPushNotificationQueue(mockQueue)
		err := svc.Enqueue(ctx, buyerID, payload)

		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}
