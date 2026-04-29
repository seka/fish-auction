package sqs

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	emailMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

type mockJobQueue struct {
	enqueueFunc func(ctx context.Context, jobType model.JobType, payload any) error
}

func (m *mockJobQueue) Enqueue(ctx context.Context, jobType model.JobType, payload any) error {
	return m.enqueueFunc(ctx, jobType, payload)
}

func (m *mockJobQueue) Dequeue(ctx context.Context, waitTimeSeconds int32) ([]*model.JobMessage, error) {
	return nil, nil
}

func (m *mockJobQueue) DeleteMessage(ctx context.Context, message *model.JobMessage) error {
	return nil
}

func TestAdminEmailService_SendAdminPasswordReset(t *testing.T) {
	ctx := context.Background()
	to := "admin@example.com"
	resetURL := "http://example.com/reset"

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

		svc := NewAdminEmailService(mockQueue)
		err := svc.SendAdminPasswordReset(ctx, to, resetURL)

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
		if p.EmailType != emailMessage.EmailTypeAdminPasswordReset {
			t.Errorf("Expected EmailType %s, got %s", emailMessage.EmailTypeAdminPasswordReset, p.EmailType)
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
		mockQueue := &mockJobQueue{
			enqueueFunc: func(_ context.Context, _ model.JobType, _ any) error {
				return expectedErr
			},
		}

		svc := NewAdminEmailService(mockQueue)
		err := svc.SendAdminPasswordReset(ctx, to, resetURL)

		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}
