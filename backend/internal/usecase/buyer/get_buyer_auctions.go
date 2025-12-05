package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type GetBuyerAuctionsUseCase interface {
	Execute(ctx context.Context, buyerID int) ([]model.Auction, error)
}

type getBuyerAuctionsUseCase struct {
	bidRepo repository.BidRepository
}

func NewGetBuyerAuctionsUseCase(bidRepo repository.BidRepository) GetBuyerAuctionsUseCase {
	return &getBuyerAuctionsUseCase{
		bidRepo: bidRepo,
	}
}

func (uc *getBuyerAuctionsUseCase) Execute(ctx context.Context, buyerID int) ([]model.Auction, error) {
	return uc.bidRepo.ListAuctionsByBuyerID(ctx, buyerID)
}
