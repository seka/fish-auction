package repository

import (
	"context"
	"time"
)

// RateLimitRepository manages per-key fixed-window counters for rate limiting.
type RateLimitRepository interface {
	// Increment atomically increments the request counter for the given IP and window,
	// and returns the updated count.
	Increment(ctx context.Context, keyPrefix string, ip string, window time.Duration) (int64, error)
}
