package push_notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

var webpushSendNotification = webpush.SendNotification

type webpushNotificationService struct {
	cfg *config.Config
}

var _ service.PushNotificationService = (*webpushNotificationService)(nil)

// NewWebpushService creates a new PushNotificationService implementation using webpush-go
func NewWebpushService(cfg *config.Config) *webpushNotificationService {
	return &webpushNotificationService{
		cfg: cfg,
	}
}

func (s *webpushNotificationService) Send(_ context.Context, sub *model.PushSubscription, payload any) error {
	message, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	pushSub := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256dh,
			Auth:   sub.Auth,
		},
	}

	resp, err := webpushSendNotification(message, pushSub, &webpush.Options{
		Subscriber:      s.cfg.VAPIDSubject,
		VAPIDPublicKey:  s.cfg.VAPIDPublicKey,
		VAPIDPrivateKey: s.cfg.VAPIDPrivateKey,
		TTL:             30, // seconds
	})

	if err != nil {
		return fmt.Errorf("failed to send push notification: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	log.Printf("Push notification sent to %s, status: %d", sub.Endpoint, resp.StatusCode)

	if resp.StatusCode == 410 || resp.StatusCode == 404 {
		return fmt.Errorf("subscription expired (status %d)", resp.StatusCode)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
