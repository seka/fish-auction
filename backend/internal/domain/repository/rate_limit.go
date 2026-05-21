package repository

import (
	"context"
	"time"
)

// RateLimitRepository manages per-IP fixed-window counters for rate limiting.
type RateLimitRepository interface {
	IncrementAdminLogin(ctx context.Context, ip string, window time.Duration) (int64, error)
	IncrementBuyerLogin(ctx context.Context, ip string, window time.Duration) (int64, error)
	IncrementAdminReset(ctx context.Context, ip string, window time.Duration) (int64, error)
	IncrementBuyerReset(ctx context.Context, ip string, window time.Duration) (int64, error)
}
