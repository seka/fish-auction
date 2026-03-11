package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

type fishermanStore struct {
	cache datastore.Cache
	ttl   time.Duration
}

// NewFishermanStore は新しい FishermanStore を作成
func NewFishermanStore(cache datastore.Cache, ttl time.Duration) datastore.FishermanCache {
	return &fishermanStore{
		cache: cache,
		ttl:   ttl,
	}
}

func (c *fishermanStore) Get(ctx context.Context, id int) (*model.Fisherman, error) {
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

func (c *fishermanStore) Set(ctx context.Context, id int, fisherman *model.Fisherman) error {
	key := fmt.Sprintf("fisherman:%d", id)
	data, err := json.Marshal(fisherman)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.ttl)
}

func (c *fishermanStore) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("fisherman:%d", id)
	return c.cache.Delete(ctx, key)
}
