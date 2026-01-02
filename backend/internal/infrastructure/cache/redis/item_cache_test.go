package cache_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	cache "github.com/seka/fish-auction/backend/internal/infrastructure/cache/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemCache_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	c := cache.NewItemCache(db, time.Hour)
	ctx := context.Background()

	t.Run("CacheHit", func(t *testing.T) {
		item := &model.AuctionItem{ID: 1, FishType: "Tuna"}
		data, _ := json.Marshal(item)
		mock.ExpectGet("item:1").SetVal(string(data))

		got, err := c.Get(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, item.ID, got.ID)
		assert.Equal(t, item.FishType, got.FishType)
	})

	t.Run("CacheMiss", func(t *testing.T) {
		mock.ExpectGet("item:1").RedisNil()

		got, err := c.Get(ctx, 1)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("RedisError", func(t *testing.T) {
		mock.ExpectGet("item:1").SetErr(errors.New("connection failed"))

		got, err := c.Get(ctx, 1)
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestItemCache_Set(t *testing.T) {
	db, mock := redismock.NewClientMock()
	ttl := time.Hour
	c := cache.NewItemCache(db, ttl)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		item := &model.AuctionItem{ID: 1, FishType: "Tuna"}
		data, _ := json.Marshal(item)
		mock.ExpectSet("item:1", data, ttl).SetVal("OK")

		err := c.Set(ctx, 1, item)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		item := &model.AuctionItem{ID: 1, FishType: "Tuna"}
		data, _ := json.Marshal(item)
		mock.ExpectSet("item:1", data, ttl).SetErr(errors.New("failed"))

		err := c.Set(ctx, 1, item)
		assert.Error(t, err)
	})
}

func TestItemCache_Delete(t *testing.T) {
	db, mock := redismock.NewClientMock()
	c := cache.NewItemCache(db, time.Hour)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectDel("item:1").SetVal(1)

		err := c.Delete(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectDel("item:1").SetErr(errors.New("failed"))

		err := c.Delete(ctx, 1)
		assert.Error(t, err)
	})
}
