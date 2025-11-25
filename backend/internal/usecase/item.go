package usecase

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ItemUseCase defines the interface for auction item-related business logic
type ItemUseCase interface {
	Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
	List(ctx context.Context, status string) ([]model.AuctionItem, error)
}

type itemInteractor struct {
	repo repository.ItemRepository
}

func NewItemInteractor(repo repository.ItemRepository) ItemUseCase {
	return &itemInteractor{repo: repo}
}

func (i *itemInteractor) Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	return i.repo.Create(ctx, item)
}

func (i *itemInteractor) List(ctx context.Context, status string) ([]model.AuctionItem, error) {
	return i.repo.List(ctx, status)
}
