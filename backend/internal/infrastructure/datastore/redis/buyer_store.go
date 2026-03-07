package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

type buyerCacheStore struct {
	cache datastore.Cache
	ttl   time.Duration
}

// NewBuyerCacheStore は新しい BuyerCache を作成
func NewBuyerCacheStore(cache datastore.Cache, ttl time.Duration) datastore.BuyerCache {
	return &buyerCacheStore{
		cache: cache,
		ttl:   ttl,
	}
}

func (c *buyerCacheStore) Get(ctx context.Context, id int) (*model.Buyer, error) {
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

func (c *buyerCacheStore) Set(ctx context.Context, id int, buyer *model.Buyer) error {
	key := fmt.Sprintf("buyer:%d", id)
	data, err := json.Marshal(buyer)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.ttl)
}

func (c *buyerCacheStore) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("buyer:%d", id)
	return c.cache.Delete(ctx, key)
}
