package registry

import (
	"fmt"

	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/infrastructure/email/mailhog"
	"github.com/seka/fish-auction/backend/internal/infrastructure/email/templates"
	pushNotification "github.com/seka/fish-auction/backend/internal/infrastructure/push_notification"
)

// Service defines the interface for creating domain services
type Service interface {
	NewPushNotificationService() service.PushNotificationService
	NewAdminEmailService() service.AdminEmailService
	NewBuyerEmailService() service.BuyerEmailService
}

type serviceRegistry struct {
	pushNotificationService service.PushNotificationService
	adminEmailService       service.AdminEmailService
	buyerEmailService       service.BuyerEmailService
}

// NewServiceRegistry creates a new Service registry
func NewServiceRegistry(cfg *config.Config) Service {
	loader, err := templates.NewTemplateLoader()
	if err != nil {
		panic(fmt.Sprintf("failed to load templates: %v", err))
	}
	adminEmailService := mailhog.NewAdminEmailService(cfg, loader)
	buyerEmailService := mailhog.NewBuyerEmailService(cfg, loader)

	pushNotificationService := pushNotification.NewWebpushService(cfg)

	return &serviceRegistry{

		pushNotificationService: pushNotificationService,
		adminEmailService:       adminEmailService,
		buyerEmailService:       buyerEmailService,
	}
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
