package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

// FishermanStore implements a cache store for fishermen in Redis.
type FishermanStore struct {
	cache datastore.Cache
	ttl   time.Duration
}

// NewFishermanStore creates a new FishermanStore instance.
func NewFishermanStore(cache datastore.Cache, ttl time.Duration) *FishermanStore {
	return &FishermanStore{
		cache: cache,
		ttl:   ttl,
	}
}

// Get provides Get related functionality.
func (c *FishermanStore) Get(ctx context.Context, id int) (*model.Fisherman, error) {
	key := fmt.Sprintf("fisherman:%d", id)
	data, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil // キャッシュミス
	}

	var fisherman model.Fisherman
	if err := json.Unmarshal(data, &fisherman); err != nil {
		return nil, err
	}
	return &fisherman, nil
}

// Set provides Set related functionality.
func (c *FishermanStore) Set(ctx context.Context, id int, fisherman *model.Fisherman) error {
	key := fmt.Sprintf("fisherman:%d", id)
	data, err := json.Marshal(fisherman)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.ttl)
}

// Delete removes a record by ID.
func (c *FishermanStore) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("fisherman:%d", id)
	return c.cache.Delete(ctx, key)
}
