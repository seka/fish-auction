package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type FishermanCache interface {
	Get(ctx context.Context, id int) (*model.Fisherman, error)
	Set(ctx context.Context, id int, fisherman *model.Fisherman, ) error
	Delete(ctx context.Context, id int) error
}

type FishermanStore interface {
	Create(ctx context.Context, name string) (*model.Fisherman, error)
	List(ctx context.Context) ([]model.Fisherman, error)
	FindByID(ctx context.Context, id int) (*model.Fisherman, error)
	Delete(ctx context.Context, id int) error
}

type FishermanCompositeStore struct {
	store FishermanStore
	cache FishermanCache
}

func NewFishermanCompositeStore(store FishermanStore, cache FishermanCache) *FishermanCompositeStore {
	return &FishermanCompositeStore{store: store, cache: cache}
}

func (s *FishermanCompositeStore) Create(ctx context.Context, name string) (*model.Fisherman, error) {
	fisherman, err := s.store.Create(ctx, name)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, fisherman.ID)
	return fisherman, nil
}

func (s *FishermanCompositeStore) List(ctx context.Context) ([]model.Fisherman, error) {
	return s.store.List(ctx)
}

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

func (s *FishermanCompositeStore) Delete(ctx context.Context, id int) error {
	if err := s.store.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}
