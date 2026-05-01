package notification

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type mockOutboxRepository struct {
	insertFunc func(ctx context.Context, jobType model.JobType, schemaVersion int, payload []byte) error
}

func (m *mockOutboxRepository) Insert(ctx context.Context, jobType model.JobType, schemaVersion int, payload []byte) error {
	return m.insertFunc(ctx, jobType, schemaVersion, payload)
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
		var capturedJobType model.JobType

		mockRepo := &mockOutboxRepository{
			insertFunc: func(_ context.Context, jt model.JobType, _ int, _ []byte) error {
				insertCalled = true
				capturedJobType = jt
				return nil
			},
		}

		uc := NewPublishNotificationUseCase(mockRepo)
		err := uc.Execute(ctx, buyerID, payload)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !insertCalled {
			t.Error("Expected OutboxRepository.Insert to be called")
		}
		if capturedJobType != model.JobTypePushNotification {
			t.Errorf("Expected JobType %s, got %s", model.JobTypePushNotification, capturedJobType)
		}
	})

	t.Run("insert error", func(t *testing.T) {
		insertErr := errors.New("insert failed")
		mockRepo := &mockOutboxRepository{
			insertFunc: func(_ context.Context, _ model.JobType, _ int, _ []byte) error {
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
