package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListItemsUseCase handles listing auction items
type ListItemsUseCase struct {
	repo repository.ItemRepository
}

// NewListItemsUseCase creates a new instance of ListItemsUseCase
func NewListItemsUseCase(repo repository.ItemRepository) *ListItemsUseCase {
	return &ListItemsUseCase{repo: repo}
}

// Execute lists auction items, optionally filtered by status
func (uc *ListItemsUseCase) Execute(ctx context.Context, status string) ([]model.AuctionItem, error) {
	return uc.repo.List(ctx, status)
}
