package cache_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFishermanCache_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	c := cache.NewFishermanCache(db, time.Hour)
	ctx := context.Background()

	t.Run("CacheHit", func(t *testing.T) {
		fisherman := &model.Fisherman{ID: 1, Name: "Fisherman A"}
		data, _ := json.Marshal(fisherman)
		mock.ExpectGet("fisherman:1").SetVal(string(data))

		got, err := c.Get(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, fisherman.ID, got.ID)
		assert.Equal(t, fisherman.Name, got.Name)
	})

	t.Run("CacheMiss", func(t *testing.T) {
		mock.ExpectGet("fisherman:1").RedisNil()

		got, err := c.Get(ctx, 1)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("RedisError", func(t *testing.T) {
		mock.ExpectGet("fisherman:1").SetErr(errors.New("connection failed"))

		got, err := c.Get(ctx, 1)
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestFishermanCache_Set(t *testing.T) {
	db, mock := redismock.NewClientMock()
	ttl := time.Hour
	c := cache.NewFishermanCache(db, ttl)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		fisherman := &model.Fisherman{ID: 1, Name: "Fisherman A"}
		data, _ := json.Marshal(fisherman)
		mock.ExpectSet("fisherman:1", data, ttl).SetVal("OK")

		err := c.Set(ctx, 1, fisherman)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		fisherman := &model.Fisherman{ID: 1, Name: "Fisherman A"}
		data, _ := json.Marshal(fisherman)
		mock.ExpectSet("fisherman:1", data, ttl).SetErr(errors.New("failed"))

		err := c.Set(ctx, 1, fisherman)
		assert.Error(t, err)
	})
}

func TestFishermanCache_Delete(t *testing.T) {
	db, mock := redismock.NewClientMock()
	c := cache.NewFishermanCache(db, time.Hour)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectDel("fisherman:1").SetVal(1)

		err := c.Delete(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectDel("fisherman:1").SetErr(errors.New("failed"))

		err := c.Delete(ctx, 1)
		assert.Error(t, err)
	})
}
