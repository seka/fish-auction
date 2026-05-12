package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrustedProxyMiddleware_Handle(t *testing.T) {
	captured := struct {
		remoteAddr string
		scheme     string
		trusted    bool
	}{}
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured.remoteAddr = r.RemoteAddr
		captured.scheme = r.URL.Scheme
		captured.trusted = IsFromTrustedProxy(r.Context())
	})

	t.Run("noop when trusted CIDRs empty", func(t *testing.T) {
		mw := NewTrustedProxyMiddleware(nil)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.5:1234"
		req.Header.Set("X-Forwarded-For", "203.0.113.10")
		req.Header.Set("X-Forwarded-Proto", "https")

		mw.Handle(next).ServeHTTP(httptest.NewRecorder(), req)

		assert.Equal(t, "10.0.0.5:1234", captured.remoteAddr)
		assert.Empty(t, captured.scheme)
		assert.False(t, captured.trusted)
	})

	t.Run("rewrites RemoteAddr and Scheme when source is trusted", func(t *testing.T) {
		mw := NewTrustedProxyMiddleware([]string{"10.0.0.0/16"})
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.1.10:443"
		req.Header.Set("X-Forwarded-For", "203.0.113.10, 10.0.1.20")
		req.Header.Set("X-Forwarded-Proto", "https")

		mw.Handle(next).ServeHTTP(httptest.NewRecorder(), req)

		assert.Equal(t, "203.0.113.10", captured.remoteAddr)
		assert.Equal(t, "https", captured.scheme)
		assert.True(t, captured.trusted)
	})

	t.Run("does not rewrite when source is not trusted", func(t *testing.T) {
		mw := NewTrustedProxyMiddleware([]string{"10.0.0.0/16"})
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		req.RemoteAddr = "203.0.113.99:1234"
		req.Header.Set("X-Forwarded-For", "1.2.3.4")
		req.Header.Set("X-Forwarded-Proto", "https")

		mw.Handle(next).ServeHTTP(httptest.NewRecorder(), req)

		assert.Equal(t, "203.0.113.99:1234", captured.remoteAddr)
		assert.Empty(t, captured.scheme)
		assert.False(t, captured.trusted)
	})

	t.Run("invalid CIDR entries are skipped", func(t *testing.T) {
		mw := NewTrustedProxyMiddleware([]string{"not-a-cidr", "10.0.0.0/16"})
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.1:80"
		req.Header.Set("X-Forwarded-For", "198.51.100.7")

		mw.Handle(next).ServeHTTP(httptest.NewRecorder(), req)

		assert.Equal(t, "198.51.100.7", captured.remoteAddr)
		assert.True(t, captured.trusted)
	})

	t.Run("walks XFF right-to-left and stops at first untrusted", func(t *testing.T) {
		mw := NewTrustedProxyMiddleware([]string{"10.0.0.0/16"})
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.1.1:80"
		// 想定: client → another-proxy(203.0.113.5) → ALB(10.0.1.20) → app
		req.Header.Set("X-Forwarded-For", "203.0.113.5, 10.0.1.20")

		mw.Handle(next).ServeHTTP(httptest.NewRecorder(), req)

		assert.Equal(t, "203.0.113.5", captured.remoteAddr)
		assert.True(t, captured.trusted)
	})
}
