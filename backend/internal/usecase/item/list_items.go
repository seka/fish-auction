package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListItemsUseCase defines the interface for listing auction items
type ListItemsUseCase interface {
	Execute(ctx context.Context) ([]model.AuctionItem, error)
}

// ListItemsUseCase handles listing auction items
type listItemsUseCase struct {
	repo repository.ItemRepository
}

var _ ListItemsUseCase = (*listItemsUseCase)(nil)

// NewListItemsUseCase creates a new instance of ListItemsUseCase
func NewListItemsUseCase(repo repository.ItemRepository) ListItemsUseCase {
	return &listItemsUseCase{repo: repo}
}

// Execute lists auction items
func (uc *listItemsUseCase) Execute(ctx context.Context) ([]model.AuctionItem, error) {
	return uc.repo.List(ctx)
}
