package datastore

import (
	"context"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// Cache defines the generic interface for cache operations
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// BuyerCache defines the interface for buyer info caching
type BuyerCache interface {
	Get(ctx context.Context, id int) (*model.Buyer, error)
	Set(ctx context.Context, id int, buyer *model.Buyer) error
	Delete(ctx context.Context, id int) error
}

// FishermanCache defines the interface for fisherman info caching
type FishermanCache interface {
	Get(ctx context.Context, id int) (*model.Fisherman, error)
	Set(ctx context.Context, id int, fisherman *model.Fisherman) error
	Delete(ctx context.Context, id int) error
}

// ItemCache defines the interface for auction item info caching
type ItemCache interface {
	Get(ctx context.Context, id int) (*model.AuctionItem, error)
	Set(ctx context.Context, id int, item *model.AuctionItem) error
	Delete(ctx context.Context, id int) error
}
