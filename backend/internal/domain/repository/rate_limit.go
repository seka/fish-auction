package repository

import (
	"context"
	"time"
)

// RateLimitRepository manages per-key fixed-window counters for rate limiting.
type RateLimitRepository interface {
	// Increment atomically increments the counter for key within the given window
	// and returns the updated count.
	Increment(ctx context.Context, key string, windowStart time.Time) (int64, error)
}
