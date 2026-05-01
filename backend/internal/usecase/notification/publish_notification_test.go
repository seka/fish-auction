package notification

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type mockOutboxRepository struct {
	insertPushFunc func(ctx context.Context, buyerID int, payload any) error
}

func (m *mockOutboxRepository) InsertEmailJob(ctx context.Context, to string, resetURL string, emailType string) error {
	return nil
}

func (m *mockOutboxRepository) InsertPushNotificationJob(ctx context.Context, buyerID int, payload any) error {
	return m.insertPushFunc(ctx, buyerID, payload)
}

func (m *mockOutboxRepository) Claim(ctx context.Context, batchSize int, instanceID string) ([]*model.OutboxMessage, error) {
	return nil, nil
}

func (m *mockOutboxRepository) MarkProcessed(ctx context.Context, ids []int64) error {
	return nil
}

func (m *mockOutboxRepository) MarkFailed(ctx context.Context, id int64, lastError string) error {
	return nil
}

func (m *mockOutboxRepository) RecoverStale(ctx context.Context, timeout time.Duration) (int64, error) {
	return 0, nil
}

func (m *mockOutboxRepository) DeleteProcessedBefore(ctx context.Context, before time.Time) (int64, error) {
	return 0, nil
}

func TestPublishNotificationUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	buyerID := 1
	payload := map[string]string{"title": "test", "body": "hello"}

	t.Run("success", func(t *testing.T) {
		insertCalled := false

		mockRepo := &mockOutboxRepository{
			insertPushFunc: func(_ context.Context, _ int, _ any) error {
				insertCalled = true
				return nil
			},
		}

		uc := NewPublishNotificationUseCase(mockRepo)
		err := uc.Execute(ctx, buyerID, payload)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !insertCalled {
			t.Error("Expected OutboxRepository.InsertPushNotificationJob to be called")
		}
	})

	t.Run("insert error", func(t *testing.T) {
		insertErr := errors.New("insert failed")
		mockRepo := &mockOutboxRepository{
			insertPushFunc: func(_ context.Context, _ int, _ any) error {
				return insertErr
			},
		}

		uc := NewPublishNotificationUseCase(mockRepo)
		err := uc.Execute(ctx, buyerID, payload)

		if !errors.Is(err, insertErr) {
			t.Errorf("Expected error %v, got %v", insertErr, err)
		}
	})
}
