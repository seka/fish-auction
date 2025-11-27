package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// FishermanCache は漁師のキャッシュインターフェース
type FishermanCache interface {
	Get(ctx context.Context, id int) (*model.Fisherman, error)
	Set(ctx context.Context, id int, fisherman *model.Fisherman) error
	Delete(ctx context.Context, id int) error
}

type fishermanCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewFishermanCache は新しいFishermanCacheを作成
func NewFishermanCache(client *redis.Client, ttl time.Duration) FishermanCache {
	return &fishermanCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *fishermanCache) Get(ctx context.Context, id int) (*model.Fisherman, error) {
	key := fmt.Sprintf("fisherman:%d", id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // キャッシュミス
	}
	if err != nil {
		return nil, err
	}

	var fisherman model.Fisherman
	if err := json.Unmarshal(data, &fisherman); err != nil {
		return nil, err
	}
	return &fisherman, nil
}

func (c *fishermanCache) Set(ctx context.Context, id int, fisherman *model.Fisherman) error {
	key := fmt.Sprintf("fisherman:%d", id)
	data, err := json.Marshal(fisherman)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *fishermanCache) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("fisherman:%d", id)
	return c.client.Del(ctx, key).Err()
}
