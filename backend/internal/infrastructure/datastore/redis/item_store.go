package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

type itemStore struct {
	cache datastore.Cache
	ttl   time.Duration
}

// NewItemStore は新しい ItemStore を作成
func NewItemStore(cache datastore.Cache, ttl time.Duration) datastore.ItemCache {
	return &itemStore{
		cache: cache,
		ttl:   ttl,
	}
}

func (c *itemStore) Get(ctx context.Context, id int) (*model.AuctionItem, error) {
	key := fmt.Sprintf("item:%d", id)
	data, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil // キャッシュミス
	}

	var item model.AuctionItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (c *itemStore) Set(ctx context.Context, id int, item *model.AuctionItem) error {
	key := fmt.Sprintf("item:%d", id)
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.ttl)
}

func (c *itemStore) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("item:%d", id)
	return c.cache.Delete(ctx, key)
}
