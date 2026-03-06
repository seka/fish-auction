package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

// FishermanCache は漁師のキャッシュインターフェース
type FishermanCache interface {
	Get(ctx context.Context, id int) (*model.Fisherman, error)
	Set(ctx context.Context, id int, fisherman *model.Fisherman) error
	Delete(ctx context.Context, id int) error
}

type fishermanCache struct {
	cache datastore.Cache
	ttl   time.Duration
}

// NewFishermanCache は新しいFishermanCacheを作成
func NewFishermanCache(cache datastore.Cache, ttl time.Duration) FishermanCache {
	return &fishermanCache{
		cache: cache,
		ttl:   ttl,
	}
}

func (c *fishermanCache) Get(ctx context.Context, id int) (*model.Fisherman, error) {
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

func (c *fishermanCache) Set(ctx context.Context, id int, fisherman *model.Fisherman) error {
	key := fmt.Sprintf("fisherman:%d", id)
	data, err := json.Marshal(fisherman)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.ttl)
}

func (c *fishermanCache) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("fisherman:%d", id)
	return c.cache.Delete(ctx, key)
}
