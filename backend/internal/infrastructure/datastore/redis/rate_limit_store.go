package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

var _ repository.RateLimitRepository = (*RateLimitStore)(nil)

const (
	keyAdminLogin = "rate:login_admin"
	keyBuyerLogin = "rate:login_buyer"
	keyAdminReset = "rate:reset_admin"
	keyBuyerReset = "rate:reset_buyer"
)

// RateLimitStore implements repository.RateLimitRepository using Redis INCR.
type RateLimitStore struct {
	client *goredis.Client
}

// NewRateLimitStore creates a new RateLimitStore.
func NewRateLimitStore(client *goredis.Client) *RateLimitStore {
	return &RateLimitStore{client: client}
}

func (s *RateLimitStore) IncrementAdminLogin(ctx context.Context, ip string, window time.Duration) (int64, error) {
	return s.increment(ctx, keyAdminLogin, ip, window)
}

func (s *RateLimitStore) IncrementBuyerLogin(ctx context.Context, ip string, window time.Duration) (int64, error) {
	return s.increment(ctx, keyBuyerLogin, ip, window)
}

func (s *RateLimitStore) IncrementAdminReset(ctx context.Context, ip string, window time.Duration) (int64, error) {
	return s.increment(ctx, keyAdminReset, ip, window)
}

func (s *RateLimitStore) IncrementBuyerReset(ctx context.Context, ip string, window time.Duration) (int64, error) {
	return s.increment(ctx, keyBuyerReset, ip, window)
}

// increment is the shared implementation for all Increment* methods.
func (s *RateLimitStore) increment(ctx context.Context, keyPrefix, ip string, window time.Duration) (int64, error) {
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
