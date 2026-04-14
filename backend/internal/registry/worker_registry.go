package registry

import (
	"context"
	"fmt"

	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/sqs"
	"github.com/seka/fish-auction/backend/internal/worker"
	"github.com/seka/fish-auction/backend/internal/worker/job"
)

// WorkerRegistry handles the initialization of the background worker.
type WorkerRegistry interface {
	NewWorker() (*worker.Worker, error)
}

type workerRegistry struct {
	cfg        *config.Config
	repoReg    Repository
	serviceReg Service
}

// NewWorkerRegistry creates a new WorkerRegistry instance.
func NewWorkerRegistry(cfg *config.Config, repoReg Repository, serviceReg Service) WorkerRegistry {
	return &workerRegistry{
		cfg:        cfg,
		repoReg:    repoReg,
		serviceReg: serviceReg,
	}
}

func (r *workerRegistry) NewWorker() (*worker.Worker, error) {
	if r.cfg.SQSQueueURL == "" {
		return nil, fmt.Errorf("SQS_QUEUE_URL is not set")
	}

	sqsClient, err := sqs.NewClient(context.Background(), r.cfg.SQSRegion, r.cfg.SQSQueueURL, r.cfg.SQSEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQS client for worker: %w", err)
	}

	dispatcher := worker.NewDispatcher()

	// Register job handlers
	pushRepo := r.repoReg.NewPushRepository()
	pushSvc := r.serviceReg.NewPushNotificationService()
	dispatcher.Register(model.JobTypePushNotification, job.NewPushNotificationHandler(pushRepo, pushSvc))

	return worker.NewWorker(sqsClient, dispatcher), nil
}
