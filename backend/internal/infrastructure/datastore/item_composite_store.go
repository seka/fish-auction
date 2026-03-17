package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

 // ItemCache defines the interface for Item caching.
type ItemCache interface {
	Get(ctx context.Context, id int) (*model.AuctionItem, error)
	Set(ctx context.Context, id int, item *model.AuctionItem) error
	Delete(ctx context.Context, id int) error
}

 // ItemStore defines the interface for Item persistence.
type ItemStore interface {
	Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
	List(ctx context.Context, status string) ([]model.AuctionItem, error)
	ListByAuction(ctx context.Context, auctionID int) ([]model.AuctionItem, error)
	FindByID(ctx context.Context, id int) (*model.AuctionItem, error)
	FindByIDWithLock(ctx context.Context, id int) (*model.AuctionItem, error)
	UpdateStatus(ctx context.Context, id int, status model.ItemStatus) error
	Update(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
	Delete(ctx context.Context, id int) error
	UpdateSortOrder(ctx context.Context, id, sortOrder int) error
	Reorder(ctx context.Context, auctionID int, ids []int) error
}

// ItemCompositeStore combines persistence and caching for auction items.
type ItemCompositeStore struct {
	store ItemStore
	cache ItemCache
}

// NewItemCompositeStore creates a new ItemCompositeStore instance.
func NewItemCompositeStore(store ItemStore, cache ItemCache) *ItemCompositeStore {
	return &ItemCompositeStore{store: store, cache: cache}
}

// Create stores a new auction item in the persistence layer and invalidates the cache.
func (s *ItemCompositeStore) Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	newItem, err := s.store.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, newItem.ID)
	return newItem, nil
}

// List returns all auction items with the given status from the persistence layer.
func (s *ItemCompositeStore) List(ctx context.Context, status string) ([]model.AuctionItem, error) {
	return s.store.List(ctx, status)
}

// ListByAuction returns all auction items for the given auction ID from the persistence layer.
func (s *ItemCompositeStore) ListByAuction(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
	return s.store.ListByAuction(ctx, auctionID)
}

// FindByID returns an auction item by its ID, checking the cache first.
func (s *ItemCompositeStore) FindByID(ctx context.Context, id int) (*model.AuctionItem, error) {
	if i, err := s.cache.Get(ctx, id); err == nil && i != nil {
		return i, nil
	}
	item, err := s.store.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, id, item)
	return item, nil
}

// FindByIDWithLock returns an auction item by its ID with a lock, bypassing the cache.
func (s *ItemCompositeStore) FindByIDWithLock(ctx context.Context, id int) (*model.AuctionItem, error) {
	// 悲観的ロックを行う場合はキャッシュをバイパスする
	item, err := s.store.FindByIDWithLock(ctx, id)
	if err != nil {
		return nil, err
	}
	// ロック取得後の値でキャッシュは更新しない（トランザクション終了後に別の場所で無効化されるため）
	return item, nil
}

// UpdateStatus updates the status of an auction item in the persistence layer and invalidates the cache.
func (s *ItemCompositeStore) UpdateStatus(ctx context.Context, id int, status model.ItemStatus) error {
	if err := s.store.UpdateStatus(ctx, id, status); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}

// Update updates an auction item in the persistence layer and invalidates the cache.
func (s *ItemCompositeStore) Update(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	updatedItem, err := s.store.Update(ctx, item)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, updatedItem.ID)
	return updatedItem, nil
}

// Delete removes an auction item by its ID from the persistence layer and the cache.
func (s *ItemCompositeStore) Delete(ctx context.Context, id int) error {
	if err := s.store.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}

// UpdateSortOrder updates the sort order of an auction item in the persistence layer and invalidates the cache.
func (s *ItemCompositeStore) UpdateSortOrder(ctx context.Context, id, sortOrder int) error {
	if err := s.store.UpdateSortOrder(ctx, id, sortOrder); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}

// Reorder reorders auction items in the persistence layer and invalidates their cache.
func (s *ItemCompositeStore) Reorder(ctx context.Context, auctionID int, ids []int) error {
	if err := s.store.Reorder(ctx, auctionID, ids); err != nil {
		return err
	}
	for _, id := range ids {
		_ = s.cache.Delete(ctx, id)
	}
	return nil
}

// InvalidateCache invalidates the cache for the given auction item ID.
func (s *ItemCompositeStore) InvalidateCache(ctx context.Context, id int) error {
	return s.cache.Delete(ctx, id)
}
