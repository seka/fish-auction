package push_notification

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/model"
)

func TestWebpushNotificationService_Send(t *testing.T) {
	cfg := &config.Config{
		VAPIDPublicKey:  "test-public-key",
		VAPIDPrivateKey: "test-private-key",
		VAPIDSubject:    "mailto:test@example.com",
	}

	subscription := &model.PushSubscription{
		Endpoint: "https://example.com/push",
		P256dh:   "p256dh-key",
		Auth:     "auth-secret",
	}

	payload := map[string]string{"message": "hello"}

	t.Run("Success", func(t *testing.T) {
		orig := webpushSendNotification
		defer func() { webpushSendNotification = orig }()

		webpushSendNotification = func(_ []byte, _ *webpush.Subscription, _ *webpush.Options) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		}

		svc := NewWebpushService(cfg)
		err := svc.Send(context.Background(), subscription, payload)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("SubscriptionExpired_410", func(t *testing.T) {
		orig := webpushSendNotification
		defer func() { webpushSendNotification = orig }()

		webpushSendNotification = func(_ []byte, _ *webpush.Subscription, _ *webpush.Options) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusGone,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		}

		svc := NewWebpushService(cfg)
		err := svc.Send(context.Background(), subscription, payload)

		var pushErr *PushNotificationError
		if !errors.As(err, &pushErr) || pushErr.StatusCode != http.StatusGone {
			t.Errorf("expected PushNotificationError with status 410, got %v", err)
		}
	})

	t.Run("SubscriptionNotFound_404", func(t *testing.T) {
		orig := webpushSendNotification
		defer func() { webpushSendNotification = orig }()

		webpushSendNotification = func(_ []byte, _ *webpush.Subscription, _ *webpush.Options) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		}

		svc := NewWebpushService(cfg)
		err := svc.Send(context.Background(), subscription, payload)

		var pushErr *PushNotificationError
		if !errors.As(err, &pushErr) || pushErr.StatusCode != http.StatusNotFound {
			t.Errorf("expected PushNotificationError with status 404, got %v", err)
		}
	})

	t.Run("UnexpectedStatus_500", func(t *testing.T) {
		orig := webpushSendNotification
		defer func() { webpushSendNotification = orig }()

		webpushSendNotification = func(_ []byte, _ *webpush.Subscription, _ *webpush.Options) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		}

		svc := NewWebpushService(cfg)
		err := svc.Send(context.Background(), subscription, payload)

		var pushErr *PushNotificationError
		if !errors.As(err, &pushErr) || pushErr.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected PushNotificationError with status 500, got %v", err)
		}
	})

	t.Run("LibraryError", func(t *testing.T) {
		orig := webpushSendNotification
		defer func() { webpushSendNotification = orig }()

		expectedErr := errors.New("network error")
		webpushSendNotification = func(_ []byte, _ *webpush.Subscription, _ *webpush.Options) (*http.Response, error) {
			return nil, expectedErr
		}

		svc := NewWebpushService(cfg)
		err := svc.Send(context.Background(), subscription, payload)

		if !errors.Is(err, expectedErr) && err.Error() != "failed to send push notification: network error" {
			t.Errorf("expected original error or wrapped error, got %v", err)
		}
	})
}
