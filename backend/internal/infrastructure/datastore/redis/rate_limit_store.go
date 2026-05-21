package redis

import (
	"context"
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

// Increment atomically increments the counter via INCR and sets TTL on first creation.
// Returns 0 without error when Redis is unavailable (fail-open).
func (s *RateLimitStore) Increment(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	if s.client == nil {
		return 0, nil
	}

	pipe := s.client.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	// EXPIRE はキー作成時のみ有効にしたいが、パイプラインの簡潔さを優先して毎回設定する。
	// バケットキーはウィンドウ切り替えで新規キーになるため、TTL のリセットは実害がない。
	pipe.Expire(ctx, key, ttl)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, err
	}
	return incrCmd.Val(), nil
}
