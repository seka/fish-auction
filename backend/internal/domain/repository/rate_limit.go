package repository

import (
	"context"
	"time"
)

// RateLimitRepository manages per-key fixed-window counters for rate limiting.
type RateLimitRepository interface {
	// Increment atomically increments the counter for key and returns the updated count.
	// ttl is applied to the key on creation to ensure automatic expiry.
	Increment(ctx context.Context, key string, ttl time.Duration) (int64, error)
}
