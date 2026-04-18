package registry

import (
	"context"
	"fmt"

	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/email/mailhog"
	"github.com/seka/fish-auction/backend/internal/infrastructure/email/templates"
	pushNotification "github.com/seka/fish-auction/backend/internal/infrastructure/push_notification"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/sqs"
)

// Service defines the interface for creating domain services
type Service interface {
	NewPushNotificationService() service.PushNotificationService
	NewAdminEmailService() service.AdminEmailService
	NewBuyerEmailService() service.BuyerEmailService
	NewClock() service.Clock
	NewJobQueue() service.JobQueue
}

type serviceRegistry struct {
	pushNotificationService service.PushNotificationService
	adminEmailService       service.AdminEmailService
	buyerEmailService       service.BuyerEmailService
	clock                   service.Clock
	jobQueue                service.JobQueue
}

// NewServiceRegistry creates a new Service registry
func NewServiceRegistry(
	emailCfg config.EmailConfig,
	webpushCfg config.WebpushConfig,
	jobQueueCfg config.QueueConfig,
) (Service, error) {
	loader, err := templates.NewTemplateLoader()
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}
	adminEmailService := mailhog.NewAdminEmailService(emailCfg, loader)
	buyerEmailService := mailhog.NewBuyerEmailService(emailCfg, loader)

	pushNotificationService := pushNotification.NewWebpushService(webpushCfg)

	region, url, endpoint := jobQueueCfg.SQSConfig()
	jobQueue, err := sqs.NewClient(context.Background(), region, url, endpoint)
	if err != nil {
		return nil, err
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

func (s *serviceRegistry) NewJobQueue() service.JobQueue {
	return s.jobQueue
}
