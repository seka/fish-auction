package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type BuyerCache interface {
	Get(ctx context.Context, id int) (*model.Buyer, error)
	Set(ctx context.Context, id int, buyer *model.Buyer) error
	Delete(ctx context.Context, id int) error
}

type BuyerStore interface {
	Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error)
	List(ctx context.Context) ([]model.Buyer, error)
	FindByID(ctx context.Context, id int) (*model.Buyer, error)
	FindByName(ctx context.Context, name string) (*model.Buyer, error)
	FindByEmail(ctx context.Context, email string) (*model.Buyer, error)
	Delete(ctx context.Context, id int) error
}

type BuyerCompositeStore struct {
	store BuyerStore
	cache BuyerCache
}

func NewBuyerCompositeStore(store BuyerStore, cache BuyerCache) *BuyerCompositeStore {
	return &BuyerCompositeStore{store: store, cache: cache}
}

func (s *BuyerCompositeStore) Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error) {
	newBuyer, err := s.store.Create(ctx, buyer)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, newBuyer.ID)
	return newBuyer, nil
}

func (s *BuyerCompositeStore) List(ctx context.Context) ([]model.Buyer, error) {
	return s.store.List(ctx)
}

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

func (s *BuyerCompositeStore) FindByName(ctx context.Context, name string) (*model.Buyer, error) {
	return s.store.FindByName(ctx, name)
}

func (s *BuyerCompositeStore) FindByEmail(ctx context.Context, email string) (*model.Buyer, error) {
	return s.store.FindByEmail(ctx, email)
}

func (s *BuyerCompositeStore) Delete(ctx context.Context, id int) error {
	if err := s.store.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}
