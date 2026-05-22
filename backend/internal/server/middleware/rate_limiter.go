package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/seka/fish-auction/backend/internal/server/util"
)

const (
	LoginRateWindow = 10 * time.Minute
	LoginRateLimit  = 20
	ResetRateWindow = 10 * time.Minute
	ResetRateLimit  = 10
)

// RateLimitFunc is the increment function signature used by RateLimiterMiddleware.
type RateLimitFunc func(ctx context.Context, ip string, window time.Duration) (int64, error)

// RateLimiterMiddleware enforces a per-IP fixed-window rate limit.
type RateLimiterMiddleware struct {
	increment RateLimitFunc
	limit     int
	window    time.Duration
}

// NewRateLimiterMiddleware creates a RateLimiterMiddleware.
func NewRateLimiterMiddleware(increment RateLimitFunc, limit int, window time.Duration) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		increment: increment,
		limit:     limit,
		window:    window,
	}
}

// Handle applies fixed-window rate limiting before passing the request to next.
func (m *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r.RemoteAddr)
		count, err := m.increment(r.Context(), ip, m.window)
		if err != nil {
			slog.Warn("rate limiter: increment error", "err", err, "ip", ip)
			next.ServeHTTP(w, r)
			return
		}

		if count > int64(m.limit) {
			windowSecs := int64(m.window.Seconds())
			retryAfter := windowSecs - (time.Now().UTC().Unix() % windowSecs)
			w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
			util.WriteError(w, http.StatusTooManyRequests, "Too Many Requests")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// extractIP strips the port from a "host:port" remote address string.
func extractIP(remoteAddr string) string {
	i := strings.LastIndex(remoteAddr, ":")
	if i < 0 {
		return remoteAddr
	}
	return remoteAddr[:i]
}
