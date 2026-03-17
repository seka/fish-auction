package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateItemUseCase defines the interface for creating auction items.
type CreateItemUseCase interface {
	// Execute creates a new auction item.
	Execute(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
}

// createItemUseCase handles the creation of auction items
type createItemUseCase struct {
	repo repository.ItemRepository
}

var _ CreateItemUseCase = (*createItemUseCase)(nil)

// NewCreateItemUseCase creates a new instance of CreateItemUseCase
func NewCreateItemUseCase(repo repository.ItemRepository) *createItemUseCase {
	return &createItemUseCase{repo: repo}
}

// Execute creates a new auction item
func (uc *createItemUseCase) Execute(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	return uc.repo.Create(ctx, item)
}
