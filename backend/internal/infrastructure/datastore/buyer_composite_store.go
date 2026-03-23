// Package datastore provides repository implementations using multiple backends.
package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// BuyerCache defines the interface for Buyer caching.
type BuyerCache interface {
	Get(ctx context.Context, id int) (*model.Buyer, error)
	Set(ctx context.Context, id int, buyer *model.Buyer) error
	Delete(ctx context.Context, id int) error
}

// BuyerStore defines the interface for Buyer persistence.
type BuyerStore interface {
	Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error)
	List(ctx context.Context) ([]model.Buyer, error)
	FindByID(ctx context.Context, id int) (*model.Buyer, error)
	FindByName(ctx context.Context, name string) (*model.Buyer, error)
	FindByEmail(ctx context.Context, email string) (*model.Buyer, error)
	Delete(ctx context.Context, id int) error
}

// BuyerCompositeStore combines persistence and caching for buyers.
type BuyerCompositeStore struct {
	store BuyerStore
	cache BuyerCache
}

// NewBuyerCompositeStore creates a new BuyerCompositeStore instance.
func NewBuyerCompositeStore(store BuyerStore, cache BuyerCache) *BuyerCompositeStore {
	return &BuyerCompositeStore{store: store, cache: cache}
}

// Create stores a new buyer in the persistence layer and invalidates the cache.
func (s *BuyerCompositeStore) Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error) {
	newBuyer, err := s.store.Create(ctx, buyer)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, newBuyer.ID)
	return newBuyer, nil
}

// List returns all buyers from the persistence layer.
func (s *BuyerCompositeStore) List(ctx context.Context) ([]model.Buyer, error) {
	return s.store.List(ctx)
}

// FindByID returns a buyer by its ID, checking the cache first.
func (s *BuyerCompositeStore) FindByID(ctx context.Context, id int) (*model.Buyer, error) {
	if b, err := s.cache.Get(ctx, id); err == nil && b != nil {
		return b, nil
	}
	buyer, err := s.store.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, id, buyer)
	return buyer, nil
}

// FindByName returns a buyer by its name from the persistence layer.
func (s *BuyerCompositeStore) FindByName(ctx context.Context, name string) (*model.Buyer, error) {
	return s.store.FindByName(ctx, name)
}

// FindByEmail returns a buyer by its email from the persistence layer.
func (s *BuyerCompositeStore) FindByEmail(ctx context.Context, email string) (*model.Buyer, error) {
	return s.store.FindByEmail(ctx, email)
}

// Delete removes a buyer by its ID from the persistence layer and the cache.
func (s *BuyerCompositeStore) Delete(ctx context.Context, id int) error {
	if err := s.store.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}
