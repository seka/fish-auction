package sqs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/mock"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/sqs"
	emailMessage "github.com/seka/fish-auction/backend/internal/worker/message"
)

func TestBuyerEmailQueue_EnqueueBuyerPasswordReset(t *testing.T) {
	ctx := context.Background()
	to := "buyer@example.com"
	resetURL := "http://example.com/reset"

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

		svc := sqs.NewBuyerEmailQueue(mockQueue)
		err := svc.EnqueueBuyerPasswordReset(ctx, to, resetURL)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !enqueueCalled {
			t.Fatal("Expected JobQueue.Enqueue to be called")
		}
		if capturedJobType != model.JobTypeEmail {
			t.Errorf("Expected jobType '%s', got '%s'", model.JobTypeEmail, capturedJobType)
		}

		p, ok := capturedPayload.(emailMessage.EmailMessage)
		if !ok {
			t.Fatalf("Captured payload is not of type EmailMessage: %T", capturedPayload)
		}
		if p.EmailType != emailMessage.EmailTypeBuyerPasswordReset {
			t.Errorf("Expected EmailType %s, got %s", emailMessage.EmailTypeBuyerPasswordReset, p.EmailType)
		}
		if p.To != to {
			t.Errorf("Expected To %s, got %s", to, p.To)
		}
		if p.ResetURL != resetURL {
			t.Errorf("Expected ResetURL %s, got %s", resetURL, p.ResetURL)
		}
	})

	t.Run("enqueue error", func(t *testing.T) {
		expectedErr := errors.New("enqueue error")
		mockQueue := &mock.MockJobQueue{
			EnqueueFunc: func(_ context.Context, _ model.JobType, _ any) error {
				return expectedErr
			},
		}

		svc := sqs.NewBuyerEmailQueue(mockQueue)
		err := svc.EnqueueBuyerPasswordReset(ctx, to, resetURL)

		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}
