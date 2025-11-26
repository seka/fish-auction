package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateItemUseCase handles the creation of auction items
type CreateItemUseCase struct {
	repo repository.ItemRepository
}

// NewCreateItemUseCase creates a new instance of CreateItemUseCase
func NewCreateItemUseCase(repo repository.ItemRepository) *CreateItemUseCase {
	return &CreateItemUseCase{repo: repo}
}

// Execute creates a new auction item
func (uc *CreateItemUseCase) Execute(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	return uc.repo.Create(ctx, item)
}
