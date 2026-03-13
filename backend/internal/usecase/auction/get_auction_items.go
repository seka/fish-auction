package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// GetAuctionItemsUseCase defines the interface for getting items in an auction
type GetAuctionItemsUseCase interface {
	Execute(ctx context.Context, auctionID int) ([]model.AuctionItem, error)
}

// getAuctionItemsUseCase handles getting items in an auction
type getAuctionItemsUseCase struct {
	itemRepo repository.ItemRepository
}

var _ GetAuctionItemsUseCase = (*getAuctionItemsUseCase)(nil)

// NewGetAuctionItemsUseCase creates a new instance of GetAuctionItemsUseCase
func NewGetAuctionItemsUseCase(itemRepo repository.ItemRepository) *getAuctionItemsUseCase {
	return &getAuctionItemsUseCase{itemRepo: itemRepo}
}

// Execute gets all items in an auction
func (uc *getAuctionItemsUseCase) Execute(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
	return uc.itemRepo.ListByAuction(ctx, auctionID)
}
