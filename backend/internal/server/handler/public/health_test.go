package public_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/server/handler/public"
)

func TestHealthHandler_Check(t *testing.T) {
	h := public.NewHealthHandler()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHealthHandler_RegisterRoutes(t *testing.T) {
	t.Run("MethodNotAllowed", func(t *testing.T) {
		h := public.NewHealthHandler()
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/health", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})
}
