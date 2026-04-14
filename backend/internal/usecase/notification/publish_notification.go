package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

// PublishNotificationUseCase defines the interface for publish notifications.
type PublishNotificationUseCase interface {
	// Execute enqueues a notification job for a buyer.
	Execute(ctx context.Context, buyerID int, payload any) error
}

type publishNotificationUseCase struct {
	jobQueue service.JobQueue
}

var _ PublishNotificationUseCase = (*publishNotificationUseCase)(nil)

// NewPublishNotificationUseCase creates a new instance of PublishNotificationUseCase.
func NewPublishNotificationUseCase(
	jobQueue service.JobQueue,
) PublishNotificationUseCase {
	return &publishNotificationUseCase{
		jobQueue: jobQueue,
	}
}

// pushNotificationJobDTO is a private DTO for marshaling the job payload to the queue.
type pushNotificationJobDTO struct {
	BuyerID int `json:"buyer_id"`
	Payload any `json:"payload"`
}

func (uc *publishNotificationUseCase) Execute(ctx context.Context, buyerID int, payload any) error {
	// Use DTO for serialization to avoid JSON tags in domain model.
	job := pushNotificationJobDTO{
		BuyerID: buyerID,
		Payload: payload,
	}

	jobBytes, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal notification job: %w", err)
	}

	// Enqueue the job for the worker to process asynchronously.
	return uc.jobQueue.Enqueue(ctx, model.JobTypePushNotification, jobBytes)
}
