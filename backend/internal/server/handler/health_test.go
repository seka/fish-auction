package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/server/handler"
)

func TestHealthHandler_Check(t *testing.T) {
	h := handler.NewHealthHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	h.Check(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "Backend is healthy!" {
		t.Errorf("expected body 'Backend is healthy!', got '%s'", w.Body.String())
	}
}

func TestHealthHandler_RegisterRoutes(t *testing.T) {
	t.Run("MethodNotAllowed", func(t *testing.T) {
		h := handler.NewHealthHandler()
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})
}
