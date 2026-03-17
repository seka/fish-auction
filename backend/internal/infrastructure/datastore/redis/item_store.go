package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

// ItemStore implements a cache store for auction items in Redis.
type ItemStore struct {
	cache datastore.Cache
	ttl   time.Duration
}

// NewItemStore creates a new ItemStore instance.
func NewItemStore(cache datastore.Cache, ttl time.Duration) *ItemStore {
	return &ItemStore{
		cache: cache,
		ttl:   ttl,
	}
}

func (c *ItemStore) Get(ctx context.Context, id int) (*model.AuctionItem, error) {
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

func (c *ItemStore) Set(ctx context.Context, id int, item *model.AuctionItem) error {
	key := fmt.Sprintf("item:%d", id)
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.ttl)
}

func (c *ItemStore) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("item:%d", id)
	return c.cache.Delete(ctx, key)
}
