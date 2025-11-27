package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// BuyerCache は買い手のキャッシュインターフェース
type BuyerCache interface {
	Get(ctx context.Context, id int) (*model.Buyer, error)
	Set(ctx context.Context, id int, buyer *model.Buyer) error
	Delete(ctx context.Context, id int) error
}

type buyerCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewBuyerCache は新しいBuyerCacheを作成
func NewBuyerCache(client *redis.Client, ttl time.Duration) BuyerCache {
	return &buyerCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *buyerCache) Get(ctx context.Context, id int) (*model.Buyer, error) {
	key := fmt.Sprintf("buyer:%d", id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // キャッシュミス
	}
	if err != nil {
		return nil, err
	}

	var buyer model.Buyer
	if err := json.Unmarshal(data, &buyer); err != nil {
		return nil, err
	}
	return &buyer, nil
}

func (c *buyerCache) Set(ctx context.Context, id int, buyer *model.Buyer) error {
	key := fmt.Sprintf("buyer:%d", id)
	data, err := json.Marshal(buyer)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *buyerCache) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("buyer:%d", id)
	return c.client.Del(ctx, key).Err()
}
