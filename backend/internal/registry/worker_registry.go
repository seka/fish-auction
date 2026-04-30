package registry

import (
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/worker/handler"
	"github.com/seka/fish-auction/backend/internal/worker"
)

// WorkerRegistry handles the initialization of the background worker.
type WorkerRegistry interface {
	NewWorker() (*worker.Worker, error)
}

type workerRegistry struct {
	repoReg    Repository
	serviceReg Service
}

// NewWorkerRegistry creates a new WorkerRegistry instance.
func NewWorkerRegistry(queueCfg config.QueueConfig, repoReg Repository, serviceReg Service) (WorkerRegistry, error) {
	return &workerRegistry{
		repoReg:    repoReg,
		serviceReg: serviceReg,
	}, nil
}

func (r *workerRegistry) NewWorker() (*worker.Worker, error) {
	// Initialize handlers
	pushRepo := r.repoReg.NewPushRepository()
	pushSvc := r.serviceReg.NewPushNotificationService()
	pushHandler := handler.NewPushNotificationHandler(pushRepo, pushSvc)

	buyerEmailSvc := r.serviceReg.NewBuyerEmailService()
	adminEmailSvc := r.serviceReg.NewAdminEmailService()
	emailHandler := handler.NewEmailHandler(buyerEmailSvc, adminEmailSvc)

	adminEmailQueue := r.serviceReg.NewAdminEmailQueue()
	buyerEmailQueue := r.serviceReg.NewBuyerEmailQueue()
	pushNotificationQueue := r.serviceReg.NewPushNotificationQueue()

	return worker.NewWorker(
		adminEmailQueue,
		buyerEmailQueue,
		pushNotificationQueue,
		worker.HandlerFunc(emailHandler.Handle),
		worker.HandlerFunc(pushHandler.Handle),
	), nil
}
