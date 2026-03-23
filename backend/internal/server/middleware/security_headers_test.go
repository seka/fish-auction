package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	middleware := NewSecurityHeadersMiddleware()

	handler := middleware.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	res := rr.Result()

	expectedHeaders := map[string]string{
		"Content-Security-Policy":           "default-src 'none'; frame-ancestors 'none'",
		"Strict-Transport-Security":         "max-age=631138519",
		"X-Content-Type-Options":            "nosniff",
		"X-Download-Options":                "noopen",
		"X-Frame-Options":                   "DENY",
		"X-Permitted-Cross-Domain-Policies": "none",
		"X-XSS-Protection":                  "0",
		"Referrer-Policy":                   "strict-origin-when-cross-origin",
	}

	for key, expectedValue := range expectedHeaders {
		if value := res.Header.Get(key); value != expectedValue {
			t.Errorf("Expected header %s to be %q, got %q", key, expectedValue, value)
		}
	}
}
