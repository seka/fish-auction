package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// UpdateItemSortOrderUseCase updates an existing record.
type UpdateItemSortOrderUseCase interface {
	Execute(ctx context.Context, id int, sortOrder int) error
}

var _ UpdateItemSortOrderUseCase = (*updateItemSortOrderUseCase)(nil)

type updateItemSortOrderUseCase struct {
	repo repository.ItemRepository
}

// NewUpdateItemSortOrderUseCase creates a new instance of UpdateItemSortOrderUseCase
func NewUpdateItemSortOrderUseCase(repo repository.ItemRepository) UpdateItemSortOrderUseCase {
	return &updateItemSortOrderUseCase{repo: repo}
}

func (uc *updateItemSortOrderUseCase) Execute(ctx context.Context, id, sortOrder int) error {
	return uc.repo.UpdateSortOrder(ctx, id, sortOrder)
}
