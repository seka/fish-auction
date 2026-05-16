package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/server/util"
)

const (
	LoginRateWindow = 10 * time.Minute
	LoginRateLimit  = 20
	ResetRateWindow = 10 * time.Minute
	ResetRateLimit  = 10
)

// RateLimiterMiddleware enforces a per-IP sliding window rate limit using Redis sorted sets.
type RateLimiterMiddleware struct {
	redisClient *redis.Client
	limit       int
	window      time.Duration
	keyPrefix   string
}

// NewRateLimiterMiddleware creates a RateLimiterMiddleware.
func NewRateLimiterMiddleware(rc *redis.Client, keyPrefix string, limit int, window time.Duration) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		redisClient: rc,
		limit:       limit,
		window:      window,
		keyPrefix:   keyPrefix,
	}
}

// Handle applies sliding-window rate limiting before passing the request to next.
func (m *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redis が設定されていない場合は制限をかけない（fail-open）。
		if m.redisClient == nil {
			next.ServeHTTP(w, r)
			return
		}

		ip := extractIP(r.RemoteAddr)
		key := fmt.Sprintf("%s:%s", m.keyPrefix, ip)
		now := time.Now()
		nowNano := float64(now.UnixNano())
		windowStart := float64(now.Add(-m.window).UnixNano())

		ctx := r.Context()

		// ウィンドウ外の古いエントリを削除してからカウントする。
		pipe := m.redisClient.Pipeline()
		pipe.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%f", windowStart))
		cardCmd := pipe.ZCard(ctx, key)
		if _, err := pipe.Exec(ctx); err != nil {
			slog.Warn("rate limiter: redis pipeline error", "err", err, "key", key)
			next.ServeHTTP(w, r)
			return
		}

		count := cardCmd.Val()
		if count >= int64(m.limit) {
			retryAfter := int(m.window.Seconds())
			w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
			util.WriteError(w, http.StatusTooManyRequests, "Too Many Requests")
			return
		}

		member := fmt.Sprintf("%d", now.UnixNano())
		if err := m.redisClient.ZAdd(ctx, key, redis.Z{Score: nowNano, Member: member}).Err(); err != nil {
			slog.Warn("rate limiter: redis zadd error", "err", err, "key", key)
			next.ServeHTTP(w, r)
			return
		}
		// キーの TTL をウィンドウ幅に合わせて更新する。
		m.redisClient.Expire(ctx, key, m.window)

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
