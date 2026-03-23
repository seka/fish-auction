package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCacheControlMiddleware(t *testing.T) {
	middleware := NewCacheControlMiddleware()
	handler := middleware.Handle(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if val := rr.Header().Get("Cache-Control"); !strings.Contains(val, "no-store") {
		t.Errorf("Expected Cache-Control to contain no-store, got %q", val)
	}
}
