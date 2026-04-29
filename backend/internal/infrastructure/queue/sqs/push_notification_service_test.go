package sqs

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	notificationMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

func TestPushNotificationService_Send(t *testing.T) {
	svc := NewPushNotificationService(&mockJobQueue{})
	err := svc.Send(context.Background(), &model.PushSubscription{}, nil)
	if err == nil {
		t.Error("Expected error from Send method, got nil")
	}
}

func TestPushNotificationService_PublishToBuyer(t *testing.T) {
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

		svc := NewPushNotificationService(mockQueue)
		err := svc.PublishToBuyer(ctx, buyerID, payload)

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
		// In a real test we'd compare payload, but it's any type so we skip deep check
	})

	t.Run("enqueue error", func(t *testing.T) {
		expectedErr := errors.New("enqueue error")
		mockQueue := &mockJobQueue{
			enqueueFunc: func(_ context.Context, _ model.JobType, _ any) error {
				return expectedErr
			},
		}

		svc := NewPushNotificationService(mockQueue)
		err := svc.PublishToBuyer(ctx, buyerID, payload)

		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}
