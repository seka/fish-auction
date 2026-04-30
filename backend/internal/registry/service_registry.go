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
	NewPushNotificationQueue() service.PushNotificationQueue
	NewAdminEmailQueue() service.AdminEmailQueue
	NewBuyerEmailQueue() service.BuyerEmailQueue
	NewAdminEmailService() service.AdminEmailService
	NewBuyerEmailService() service.BuyerEmailService
	NewClock() service.Clock
}

type serviceRegistry struct {
	pushNotificationService service.PushNotificationService
	pushNotificationQueue   service.PushNotificationQueue
	adminEmailQueue         service.AdminEmailQueue
	buyerEmailQueue         service.BuyerEmailQueue
	adminEmailService       service.AdminEmailService
	buyerEmailService       service.BuyerEmailService
	clock                   service.Clock
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

	// Initialize queue clients whenever a job queue is available.
	var adminEmailService service.AdminEmailService
	var buyerEmailService service.BuyerEmailService
	var pushNotificationService service.PushNotificationService

	var pushNotificationQueue service.PushNotificationQueue
	var adminEmailQueue service.AdminEmailQueue
	var buyerEmailQueue service.BuyerEmailQueue
	if jobQueue != nil {
		adminEmailQueue = sqs.NewAdminEmailQueue(jobQueue)
		buyerEmailQueue = sqs.NewBuyerEmailQueue(jobQueue)
		pushNotificationQueue = sqs.NewPushNotificationQueue(jobQueue)
	}

	if isWorker {
		if jobQueue == nil {
			return nil, fmt.Errorf("jobQueue is required for worker process")
		}
		loader, err := templates.NewTemplateLoader()
		if err != nil {
			return nil, fmt.Errorf("failed to load templates: %w", err)
		}
		adminEmailService = mailhog.NewAdminEmailService(emailCfg, loader)
		buyerEmailService = mailhog.NewBuyerEmailService(emailCfg, loader)
		pushNotificationService = pushNotification.NewWebpushService(webpushCfg)
	} else if jobQueue == nil {
		return nil, fmt.Errorf("jobQueue is required for non-worker process")
	}

	return &serviceRegistry{
		pushNotificationService: pushNotificationService,
		pushNotificationQueue:   pushNotificationQueue,
		adminEmailQueue:         adminEmailQueue,
		buyerEmailQueue:         buyerEmailQueue,
		adminEmailService:       adminEmailService,
		buyerEmailService:       buyerEmailService,
		clock:                   service.NewRealClock(),
	}, nil
}

func (s *serviceRegistry) NewPushNotificationService() service.PushNotificationService {
	return s.pushNotificationService
}

func (s *serviceRegistry) NewPushNotificationQueue() service.PushNotificationQueue {
	return s.pushNotificationQueue
}

func (s *serviceRegistry) NewAdminEmailQueue() service.AdminEmailQueue {
	return s.adminEmailQueue
}

func (s *serviceRegistry) NewBuyerEmailQueue() service.BuyerEmailQueue {
	return s.buyerEmailQueue
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
