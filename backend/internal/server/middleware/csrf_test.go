package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFMiddleware(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000"}
	middleware := NewCSRFMiddleware(allowedOrigins)

	handler := middleware.Handle(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name           string
		method         string
		origin         string
		fetchSite      string
		expectedStatus int
	}{
		{
			name:           "GET is always allowed",
			method:         http.MethodGet,
			origin:         "http://evil.com",
			fetchSite:      "cross-site",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST with valid origin is allowed",
			method:         http.MethodPost,
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST with invalid origin is forbidden",
			method:         http.MethodPost,
			origin:         "http://evil.com",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "POST with same-origin Sec-Fetch-Site is allowed",
			method:         http.MethodPost,
			fetchSite:      "same-origin",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST with same-site Sec-Fetch-Site is allowed",
			method:         http.MethodPost,
			fetchSite:      "same-site",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST with cross-site Sec-Fetch-Site is forbidden",
			method:         http.MethodPost,
			fetchSite:      "cross-site",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "POST with no headers is allowed (compatibility)",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(context.Background(), tt.method, "/", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.fetchSite != "" {
				req.Header.Set("Sec-Fetch-Site", tt.fetchSite)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %v, got %v", tt.expectedStatus, status)
			}
		})
	}
}
