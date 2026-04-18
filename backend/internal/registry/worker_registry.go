package registry

import (
	"context"

	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/sqs"
	"github.com/seka/fish-auction/backend/internal/worker"
	"github.com/seka/fish-auction/backend/internal/worker/job"
)

// WorkerRegistry handles the initialization of the background worker.
type WorkerRegistry interface {
	NewWorker() (*worker.Worker, error)
}

type workerRegistry struct {
	queue      service.JobQueue
	repoReg    Repository
	serviceReg Service
}

// NewWorkerRegistry creates a new WorkerRegistry instance.
func NewWorkerRegistry(queueCfg config.QueueConfig, repoReg Repository, serviceReg Service) (WorkerRegistry, error) {
	region, url, endpoint := queueCfg.SQSConfig()
	queue, err := sqs.NewClient(context.Background(), region, url, endpoint)
	if err != nil {
		return nil, err
	}

	return &workerRegistry{
		queue:      queue,
		repoReg:    repoReg,
		serviceReg: serviceReg,
	}, nil
}

func (r *workerRegistry) NewWorker() (*worker.Worker, error) {
	dispatcher := worker.NewDispatcher()

	// Register job handlers
	pushRepo := r.repoReg.NewPushRepository()
	pushSvc := r.serviceReg.NewPushNotificationService()
	dispatcher.Register(model.JobTypePushNotification, job.NewPushNotificationHandler(pushRepo, pushSvc))

	return worker.NewWorker(r.queue, dispatcher), nil
}
