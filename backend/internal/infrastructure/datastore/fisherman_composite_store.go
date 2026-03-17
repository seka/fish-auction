package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

 // FishermanCache defines the interface for Fisherman caching.
type FishermanCache interface {
	Get(ctx context.Context, id int) (*model.Fisherman, error)
	Set(ctx context.Context, id int, fisherman *model.Fisherman, ) error
	Delete(ctx context.Context, id int) error
}

 // FishermanStore defines the interface for Fisherman persistence.
type FishermanStore interface {
	Create(ctx context.Context, name string) (*model.Fisherman, error)
	List(ctx context.Context) ([]model.Fisherman, error)
	FindByID(ctx context.Context, id int) (*model.Fisherman, error)
	Delete(ctx context.Context, id int) error
}

// FishermanCompositeStore combines persistence and caching for fishermen.
type FishermanCompositeStore struct {
	store FishermanStore
	cache FishermanCache
}

// NewFishermanCompositeStore creates a new FishermanCompositeStore instance.
func NewFishermanCompositeStore(store FishermanStore, cache FishermanCache) *FishermanCompositeStore {
	return &FishermanCompositeStore{store: store, cache: cache}
}

// Create stores a new fisherman in the persistence layer and invalidates the cache.
func (s *FishermanCompositeStore) Create(ctx context.Context, name string) (*model.Fisherman, error) {
	fisherman, err := s.store.Create(ctx, name)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, fisherman.ID)
	return fisherman, nil
}

// List returns all fishermen from the persistence layer.
func (s *FishermanCompositeStore) List(ctx context.Context) ([]model.Fisherman, error) {
	return s.store.List(ctx)
}

// FindByID returns a fisherman by its ID, checking the cache first.
func (s *FishermanCompositeStore) FindByID(ctx context.Context, id int) (*model.Fisherman, error) {
	if f, err := s.cache.Get(ctx, id); err == nil && f != nil {
		return f, nil
	}
	fisherman, err := s.store.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, id, fisherman)
	return fisherman, nil
}

// Delete removes a fisherman by its ID from the persistence layer and the cache.
func (s *FishermanCompositeStore) Delete(ctx context.Context, id int) error {
	if err := s.store.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}
