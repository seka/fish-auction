package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type ItemCache interface {
	Get(ctx context.Context, id int) (*model.AuctionItem, error)
	Set(ctx context.Context, id int, item *model.AuctionItem) error
	Delete(ctx context.Context, id int) error
}

type ItemCompositeStore struct {
	db    repository.ItemRepository
	cache ItemCache
}

func NewItemCompositeStore(db repository.ItemRepository, cache ItemCache) *ItemCompositeStore {
	return &ItemCompositeStore{db: db, cache: cache}
}

func (s *ItemCompositeStore) Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	newItem, err := s.db.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, newItem.ID)
	return newItem, nil
}

func (s *ItemCompositeStore) List(ctx context.Context, status string) ([]model.AuctionItem, error) {
	return s.db.List(ctx, status)
}

func (s *ItemCompositeStore) ListByAuction(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
	return s.db.ListByAuction(ctx, auctionID)
}

func (s *ItemCompositeStore) FindByID(ctx context.Context, id int) (*model.AuctionItem, error) {
	if i, err := s.cache.Get(ctx, id); err == nil && i != nil {
		return i, nil
	}
	item, err := s.db.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, id, item)
	return item, nil
}

func (s *ItemCompositeStore) UpdateStatus(ctx context.Context, id int, status model.ItemStatus) error {
	if err := s.db.UpdateStatus(ctx, id, status); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}

func (s *ItemCompositeStore) Update(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	updatedItem, err := s.db.Update(ctx, item)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, updatedItem.ID)
	return updatedItem, nil
}

func (s *ItemCompositeStore) Delete(ctx context.Context, id int) error {
	if err := s.db.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}

func (s *ItemCompositeStore) UpdateSortOrder(ctx context.Context, id int, sortOrder int) error {
	if err := s.db.UpdateSortOrder(ctx, id, sortOrder); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}

func (s *ItemCompositeStore) Reorder(ctx context.Context, auctionID int, ids []int) error {
	if err := s.db.Reorder(ctx, auctionID, ids); err != nil {
		return err
	}
	for _, id := range ids {
		_ = s.cache.Delete(ctx, id)
	}
	return nil
}

func (s *ItemCompositeStore) InvalidateCache(ctx context.Context, id int) error {
	return s.cache.Delete(ctx, id)
}
