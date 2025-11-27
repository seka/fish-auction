package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// ItemCache はアイテムのキャッシュインターフェース
type ItemCache interface {
	Get(ctx context.Context, id int) (*model.AuctionItem, error)
	Set(ctx context.Context, id int, item *model.AuctionItem) error
	Delete(ctx context.Context, id int) error
}

type itemCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewItemCache は新しいItemCacheを作成
func NewItemCache(client *redis.Client, ttl time.Duration) ItemCache {
	return &itemCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *itemCache) Get(ctx context.Context, id int) (*model.AuctionItem, error) {
	key := fmt.Sprintf("item:%d", id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // キャッシュミス
	}
	if err != nil {
		return nil, err
	}

	var item model.AuctionItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (c *itemCache) Set(ctx context.Context, id int, item *model.AuctionItem) error {
	key := fmt.Sprintf("item:%d", id)
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *itemCache) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("item:%d", id)
	return c.client.Del(ctx, key).Err()
}
