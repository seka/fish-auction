package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// PushNotificationUseCase defines the interface for push notifications
type PushNotificationUseCase interface {
	Subscribe(ctx context.Context, buyerID int, sub *model.PushSubscription) error
	SendNotification(ctx context.Context, buyerID int, payload interface{}) error
}

type pushNotificationUseCase struct {
	repo            repository.PushRepository
	vapidPublicKey  string
	vapidPrivateKey string
	vapidSubject    string
}

// NewPushNotificationUseCase creates a new instance
func NewPushNotificationUseCase(
	repo repository.PushRepository,
	publicKey, privateKey, subject string,
) PushNotificationUseCase {
	return &pushNotificationUseCase{
		repo:            repo,
		vapidPublicKey:  publicKey,
		vapidPrivateKey: privateKey,
		vapidSubject:    subject,
	}
}

func (uc *pushNotificationUseCase) Subscribe(ctx context.Context, buyerID int, sub *model.PushSubscription) error {
	sub.BuyerID = buyerID
	return uc.repo.SaveSubscription(ctx, sub)
}

func (uc *pushNotificationUseCase) SendNotification(ctx context.Context, buyerID int, payload interface{}) error {
	subs, err := uc.repo.GetSubscriptionsByBuyerID(ctx, buyerID)
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		return nil
	}

	message, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send to all subscriptions for the user
	// In production, this should probably be done asynchronously via a queue
	for _, sub := range subs {
		log.Printf("Attempting to send push notification to buyer %d, endpoint: %s", buyerID, sub.Endpoint)
		s := &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256dh,
				Auth:   sub.Auth,
			},
		}

		resp, err := webpush.SendNotification(message, s, &webpush.Options{
			Subscriber:      uc.vapidSubject,
			VAPIDPublicKey:  uc.vapidPublicKey,
			VAPIDPrivateKey: uc.vapidPrivateKey,
			TTL:             30, // seconds
		})

		if err != nil {
			log.Printf("Failed to send push notification to %s: %v", sub.Endpoint, err)
			continue
		}
		log.Printf("Push notification sent to %s, status: %d", sub.Endpoint, resp.StatusCode)
		defer resp.Body.Close()

		if resp.StatusCode == 410 || resp.StatusCode == 404 {
			log.Printf("Subscription expired for %s, deleting...", sub.Endpoint)
			_ = uc.repo.DeleteSubscription(ctx, sub.Endpoint)
		}
	}

	return nil
}
