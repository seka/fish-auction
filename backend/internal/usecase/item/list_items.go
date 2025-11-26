package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListItemsUseCase defines the interface for listing auction items
type ListItemsUseCase interface {
	Execute(ctx context.Context, status string) ([]model.AuctionItem, error)
}

// listItemsUseCase handles listing auction items
type listItemsUseCase struct {
	repo repository.ItemRepository
}

// NewListItemsUseCase creates a new instance of ListItemsUseCase
func NewListItemsUseCase(repo repository.ItemRepository) ListItemsUseCase {
	return &listItemsUseCase{repo: repo}
}

// Execute lists auction items by status
func (uc *listItemsUseCase) Execute(ctx context.Context, status string) ([]model.AuctionItem, error) {
	return uc.repo.List(ctx, status)
}
