package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type GetBuyerPurchasesUseCase interface {
	Execute(ctx context.Context, buyerID int) ([]model.Purchase, error)
}

type getBuyerPurchasesUseCase struct {
	bidRepo repository.BidRepository
}

var _ GetBuyerPurchasesUseCase = (*getBuyerPurchasesUseCase)(nil)

func NewGetBuyerPurchasesUseCase(bidRepo repository.BidRepository) *getBuyerPurchasesUseCase {
	return &getBuyerPurchasesUseCase{
		bidRepo: bidRepo,
	}
}

func (uc *getBuyerPurchasesUseCase) Execute(ctx context.Context, buyerID int) ([]model.Purchase, error) {
	return uc.bidRepo.ListPurchasesByBuyerID(ctx, buyerID)
}
