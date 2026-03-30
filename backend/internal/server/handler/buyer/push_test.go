package buyer_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/server/handler/buyer"
	"github.com/seka/fish-auction/backend/internal/server/handler/buyer/request"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestPushHandler_Subscribe(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPushUC := &mock.MockPushNotificationUseCase{
			SubscribeFunc: func(_ context.Context, buyerID int, sub *model.PushSubscription) error {
				if buyerID != 1 {
					t.Errorf("expected buyerID 1, got %d", buyerID)
				}
				if sub.Endpoint != "https://example.com/push" {
					t.Errorf("expected endpoint, got %s", sub.Endpoint)
				}
				return nil
			},
		}
		mockReg := &mock.MockRegistry{PushNotificationUC: mockPushUC}
		h := buyer.NewPushHandler(mockReg)

		reqBody := request.SubscribePush{
			Endpoint: "https://example.com/push",
		}
		reqBody.Keys.P256dh = "p256dh"
		reqBody.Keys.Auth = "auth"
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/push/subscribe", bytes.NewReader(body))
		ctx := middleware.WithBuyerID(req.Context(), 1)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.Subscribe(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Unauthorized_NoContext", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := buyer.NewPushHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/push/subscribe", bytes.NewReader([]byte("{}")))
		w := httptest.NewRecorder()

		h.Subscribe(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}
