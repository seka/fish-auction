package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type fishermanStore struct {
	db    repository.FishermanRepository
	cache FishermanCache
}

func NewFishermanRepository(db repository.FishermanRepository, cache FishermanCache) repository.FishermanRepository {
	return &fishermanStore{db: db, cache: cache}
}

func (s *fishermanStore) Create(ctx context.Context, name string) (*model.Fisherman, error) {
	fisherman, err := s.db.Create(ctx, name)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, fisherman.ID)
	return fisherman, nil
}

func (s *fishermanStore) List(ctx context.Context) ([]model.Fisherman, error) {
	return s.db.List(ctx)
}

func (s *fishermanStore) FindByID(ctx context.Context, id int) (*model.Fisherman, error) {
	if f, err := s.cache.Get(ctx, id); err == nil && f != nil {
		return f, nil
	}
	fisherman, err := s.db.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, id, fisherman)
	return fisherman, nil
}

func (s *fishermanStore) Delete(ctx context.Context, id int) error {
	if err := s.db.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}

func (s *fishermanStore) InvalidateCache(ctx context.Context, id int) error {
	return s.cache.Delete(ctx, id)
}
