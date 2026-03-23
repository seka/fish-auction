package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// GetBuyerPurchasesUseCase provides GetBuyerPurchasesUseCase related functionality.
type GetBuyerPurchasesUseCase interface {
	Execute(ctx context.Context, buyerID int) ([]model.Purchase, error)
}

type getBuyerPurchasesUseCase struct {
	bidRepo repository.BidRepository
}

var _ GetBuyerPurchasesUseCase = (*getBuyerPurchasesUseCase)(nil)

// NewGetBuyerPurchasesUseCase creates a new GetBuyerPurchasesUseCase instance.
func NewGetBuyerPurchasesUseCase(bidRepo repository.BidRepository) GetBuyerPurchasesUseCase {
	return &getBuyerPurchasesUseCase{
		bidRepo: bidRepo,
	}
}

func (uc *getBuyerPurchasesUseCase) Execute(ctx context.Context, buyerID int) ([]model.Purchase, error) {
	return uc.bidRepo.ListPurchasesByBuyerID(ctx, buyerID)
}
