package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

// BuyerStore implements a cache store for buyers in Redis.
type BuyerStore struct {
	cache datastore.Cache
	ttl   time.Duration
}

// NewBuyerStore creates a new BuyerStore instance.
func NewBuyerStore(cache datastore.Cache, ttl time.Duration) *BuyerStore {
	return &BuyerStore{
		cache: cache,
		ttl:   ttl,
	}
}

// Get provides Get related functionality.
func (c *BuyerStore) Get(ctx context.Context, id int) (*model.Buyer, error) {
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

// Set provides Set related functionality.
func (c *BuyerStore) Set(ctx context.Context, id int, buyer *model.Buyer) error {
	key := fmt.Sprintf("buyer:%d", id)
	data, err := json.Marshal(buyer)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.ttl)
}

// Delete removes a record by ID.
func (c *BuyerStore) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("buyer:%d", id)
	return c.cache.Delete(ctx, key)
}
