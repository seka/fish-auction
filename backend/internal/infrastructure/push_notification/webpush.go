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
	"github.com/seka/fish-auction/backend/internal/infrastructure/push_notification/errors"
)

// NOTE: テストでモックと差し替える際に利用している
var webpushSendNotification = webpush.SendNotification

// WebpushNotificationService provides WebpushNotificationService related functionality.
type WebpushNotificationService struct {
	cfg *config.Config
}

var _ service.PushNotificationService = (*WebpushNotificationService)(nil)

// NewWebpushService creates a new PushNotificationService implementation using webpush-go
func NewWebpushService(cfg *config.Config) *WebpushNotificationService {
	return &WebpushNotificationService{
		cfg: cfg,
	}
}

// Send provides Send related functionality.
func (s *WebpushNotificationService) Send(_ context.Context, sub *model.PushSubscription, payload any) error {
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
		return errors.HandleError(err, sub.Endpoint)
	}
	defer func() { _ = resp.Body.Close() }()

	log.Printf("Push notification sent to %s, status: %d", sub.Endpoint, resp.StatusCode)

	if resp.StatusCode == 410 || resp.StatusCode == 404 {
		return errors.HandleError(&errors.PushNotificationError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("subscription expired (status %d)", resp.StatusCode),
		}, sub.Endpoint)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.HandleError(&errors.PushNotificationError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
		}, sub.Endpoint)
	}

	return nil
}
