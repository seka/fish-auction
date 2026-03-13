package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type ReorderItemsUseCase interface {
	Execute(ctx context.Context, auctionID int, ids []int) error
}

var _ ReorderItemsUseCase = (*reorderItemsUseCase)(nil)

type reorderItemsUseCase struct {
	itemRepo repository.ItemRepository
}

// NewReorderItemsUseCase creates a new instance of ReorderItemsUseCase
func NewReorderItemsUseCase(itemRepo repository.ItemRepository) *reorderItemsUseCase {
	return &reorderItemsUseCase{itemRepo: itemRepo}
}

func (u *reorderItemsUseCase) Execute(ctx context.Context, auctionID int, ids []int) error {
	return u.itemRepo.Reorder(ctx, auctionID, ids)
}
