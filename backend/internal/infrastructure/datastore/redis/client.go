package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

type cache struct {
	client *redis.Client
}

// NewClient creates a new Cache implementation using Redis
func NewClient(client *redis.Client) datastore.Cache {
	return &cache{client: client}
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
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
