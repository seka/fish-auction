package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

// Client implements datastore.Cache using Redis.
type Client struct {
	client *goredis.Client
}

var _ datastore.Cache = (*Client)(nil)

// NewClient creates a new Client instance.
func NewClient(client *goredis.Client) *Client {
	return &Client{client: client}
}

// Get provides Get related functionality.
func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == goredis.Nil {
		return nil, nil // キャッシュミス
	}
	return data, err
}

// Set provides Set related functionality.
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a record by ID.
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	return c.client.Close()
}
