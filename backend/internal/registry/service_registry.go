package registry

import (
	"context"
	"fmt"

	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/email/mailhog"
	"github.com/seka/fish-auction/backend/internal/infrastructure/email/templates"
	pushNotification "github.com/seka/fish-auction/backend/internal/infrastructure/push_notification"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/sqs"
)

// Service defines the interface for creating domain services
type Service interface {
	NewPushNotificationService() service.PushNotificationService
	NewAdminEmailService() service.AdminEmailService
	NewBuyerEmailService() service.BuyerEmailService
	NewClock() service.Clock
	NewJobQueue() queue.JobQueue
}

type serviceRegistry struct {
	pushNotificationService service.PushNotificationService
	adminEmailService       service.AdminEmailService
	buyerEmailService       service.BuyerEmailService
	clock                   service.Clock
	jobQueue                queue.JobQueue
}

// NewServiceRegistry creates a new Service registry
func NewServiceRegistry(
	emailCfg config.EmailConfig,
	webpushCfg config.WebpushConfig,
	jobQueueCfg config.QueueConfig,
	isWorker bool,
) (Service, error) {
	var jobQueue queue.JobQueue
	if jobQueueCfg != config.NoQueueConfig {
		region, url, endpoint := jobQueueCfg.SQSConfig()
		var err error
		jobQueue, err = sqs.NewClient(context.Background(), region, url, endpoint)
		if err != nil {
			return nil, err
		}
	}

	// Use SQS-backed email/push services when a job queue is available (API server),
	// otherwise use SMTP/Webpush directly (Worker).
	var adminEmailService service.AdminEmailService
	var buyerEmailService service.BuyerEmailService
	var pushNotificationService service.PushNotificationService

	if !isWorker {
		if jobQueue == nil {
			return nil, fmt.Errorf("jobQueue is required for non-worker process")
		}
		adminEmailService = sqs.NewAdminEmailService(jobQueue)
		buyerEmailService = sqs.NewBuyerEmailService(jobQueue)
		pushNotificationService = sqs.NewPushNotificationService(jobQueue)
	} else {
		loader, err := templates.NewTemplateLoader()
		if err != nil {
			return nil, fmt.Errorf("failed to load templates: %w", err)
		}
		adminEmailService = mailhog.NewAdminEmailService(emailCfg, loader)
		buyerEmailService = mailhog.NewBuyerEmailService(emailCfg, loader)
		pushNotificationService = pushNotification.NewWebpushService(webpushCfg)
	}

	return &serviceRegistry{
		pushNotificationService: pushNotificationService,
		adminEmailService:       adminEmailService,
		buyerEmailService:       buyerEmailService,
		clock:                   service.NewRealClock(),
		jobQueue:                jobQueue,
	}, nil
}

func (s *serviceRegistry) NewPushNotificationService() service.PushNotificationService {
	return s.pushNotificationService
}

func (s *serviceRegistry) NewAdminEmailService() service.AdminEmailService {
	return s.adminEmailService
}

func (s *serviceRegistry) NewBuyerEmailService() service.BuyerEmailService {
	return s.buyerEmailService
}

func (s *serviceRegistry) NewClock() service.Clock {
	return s.clock
}

func (s *serviceRegistry) NewJobQueue() queue.JobQueue {
	return s.jobQueue
}
