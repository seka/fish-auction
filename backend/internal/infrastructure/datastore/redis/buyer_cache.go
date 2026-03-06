package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

// BuyerCache は買い手のキャッシュインターフェース
type BuyerCache interface {
	Get(ctx context.Context, id int) (*model.Buyer, error)
	Set(ctx context.Context, id int, buyer *model.Buyer) error
	Delete(ctx context.Context, id int) error
}

type buyerCache struct {
	cache datastore.Cache
	ttl   time.Duration
}

// NewBuyerCache は新しいBuyerCacheを作成
func NewBuyerCache(cache datastore.Cache, ttl time.Duration) BuyerCache {
	return &buyerCache{
		cache: cache,
		ttl:   ttl,
	}
}

func (c *buyerCache) Get(ctx context.Context, id int) (*model.Buyer, error) {
	key := fmt.Sprintf("buyer:%d", id)
	data, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil // キャッシュミス
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
	return c.cache.Set(ctx, key, data, c.ttl)
}

func (c *buyerCache) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("buyer:%d", id)
	return c.cache.Delete(ctx, key)
}
