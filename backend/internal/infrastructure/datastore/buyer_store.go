package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type BuyerCache interface {
	Get(ctx context.Context, id int) (*model.Buyer, error)
	Set(ctx context.Context, id int, buyer *model.Buyer) error
	Delete(ctx context.Context, id int) error
}

type buyerStore struct {
	db    repository.BuyerRepository
	cache BuyerCache
}

func NewBuyerCompositeStore(db repository.BuyerRepository, cache BuyerCache) repository.BuyerRepository {
	return &buyerStore{db: db, cache: cache}
}

func (s *buyerStore) Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error) {
	newBuyer, err := s.db.Create(ctx, buyer)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, newBuyer.ID)
	return newBuyer, nil
}

func (s *buyerStore) List(ctx context.Context) ([]model.Buyer, error) {
	return s.db.List(ctx)
}

func (s *buyerStore) FindByID(ctx context.Context, id int) (*model.Buyer, error) {
	if b, err := s.cache.Get(ctx, id); err == nil && b != nil {
		return b, nil
	}
	buyer, err := s.db.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, id, buyer)
	return buyer, nil
}

func (s *buyerStore) FindByName(ctx context.Context, name string) (*model.Buyer, error) {
	return s.db.FindByName(ctx, name)
}

func (s *buyerStore) FindByEmail(ctx context.Context, email string) (*model.Buyer, error) {
	return s.db.FindByEmail(ctx, email)
}

func (s *buyerStore) Delete(ctx context.Context, id int) error {
	if err := s.db.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}
