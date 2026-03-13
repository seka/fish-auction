package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type UpdateItemUseCase interface {
	Execute(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
}

var _ UpdateItemUseCase = (*updateItemUseCase)(nil)

type updateItemUseCase struct {
	repo repository.ItemRepository
}

// NewUpdateItemUseCase creates a new instance of UpdateItemUseCase
func NewUpdateItemUseCase(repo repository.ItemRepository) *updateItemUseCase {
	return &updateItemUseCase{repo: repo}
}

func (uc *updateItemUseCase) Execute(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	return uc.repo.Update(ctx, item)
}
