package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

var _ repository.RateLimitRepository = (*RateLimitStore)(nil)

// RateLimitStore implements repository.RateLimitRepository using Redis INCR.
type RateLimitStore struct {
	client *goredis.Client
}

// NewRateLimitStore creates a new RateLimitStore.
func NewRateLimitStore(client *goredis.Client) *RateLimitStore {
	return &RateLimitStore{client: client}
}

// Increment atomically increments the request counter for the given IP and window.
// The key is constructed internally as "{keyPrefix}:{bucket}:{ip}".
// Returns 0 without error when Redis is unavailable (fail-open).
func (s *RateLimitStore) Increment(ctx context.Context, keyPrefix string, ip string, window time.Duration) (int64, error) {
	if s.client == nil {
		return 0, nil
	}

	bucket := time.Now().UTC().Unix() / int64(window.Seconds())
	key := fmt.Sprintf("%s:%d:%s", keyPrefix, bucket, ip)

	pipe := s.client.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	// バケットキーはウィンドウ切り替えで新規キーになるため、TTL のリセットは実害がない。
	pipe.Expire(ctx, key, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, err
	}
	return incrCmd.Val(), nil
}
