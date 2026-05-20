package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/util"
)

const (
	LoginRateWindow = 10 * time.Minute
	LoginRateLimit  = 20
	ResetRateWindow = 10 * time.Minute
	ResetRateLimit  = 10
)

// RateLimiterMiddleware enforces a per-IP fixed-window rate limit backed by PostgreSQL.
type RateLimiterMiddleware struct {
	repo      repository.RateLimitRepository
	limit     int
	window    time.Duration
	keyPrefix string
}

// NewRateLimiterMiddleware creates a RateLimiterMiddleware.
func NewRateLimiterMiddleware(repo repository.RateLimitRepository, keyPrefix string, limit int, window time.Duration) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		repo:      repo,
		limit:     limit,
		window:    window,
		keyPrefix: keyPrefix,
	}
}

// Handle applies fixed-window rate limiting before passing the request to next.
func (m *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r.RemoteAddr)
		windowStart := time.Now().UTC().Truncate(m.window)
		key := fmt.Sprintf("%s:%s", m.keyPrefix, ip)

		count, err := m.repo.Increment(r.Context(), key, windowStart)
		if err != nil {
			slog.Warn("rate limiter: db error", "err", err, "key", key)
			next.ServeHTTP(w, r)
			return
		}

		if count > int64(m.limit) {
			retryAfter := int(m.window.Seconds())
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
