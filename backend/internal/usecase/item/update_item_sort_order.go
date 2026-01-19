package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type UpdateItemSortOrderUseCase interface {
	Execute(ctx context.Context, id int, sortOrder int) error
}

type updateItemSortOrderUseCase struct {
	repo repository.ItemRepository
}

func NewUpdateItemSortOrderUseCase(repo repository.ItemRepository) UpdateItemSortOrderUseCase {
	return &updateItemSortOrderUseCase{repo: repo}
}

func (uc *updateItemSortOrderUseCase) Execute(ctx context.Context, id int, sortOrder int) error {
	return uc.repo.UpdateSortOrder(ctx, id, sortOrder)
}
