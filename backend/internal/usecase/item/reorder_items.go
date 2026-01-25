package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type ReorderItemsUseCase interface {
	Execute(ctx context.Context, auctionID int, ids []int) error
}

type reorderItemsUseCase struct {
	itemRepo repository.ItemRepository
}

func NewReorderItemsUseCase(itemRepo repository.ItemRepository) ReorderItemsUseCase {
	return &reorderItemsUseCase{
		itemRepo: itemRepo,
	}
}

func (u *reorderItemsUseCase) Execute(ctx context.Context, auctionID int, ids []int) error {
	return u.itemRepo.Reorder(ctx, auctionID, ids)
}
