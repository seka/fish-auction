package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

type cache struct {
	client *goredis.Client
}

var _ datastore.Cache = (*cache)(nil)

// NewClient creates a new Cache implementation using Redis
func NewClient(client *goredis.Client) *cache {
	return &cache{client: client}
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == goredis.Nil {
		return nil, nil // キャッシュミス
	}
	return data, err
}

func (c *cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
