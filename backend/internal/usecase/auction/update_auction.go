package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// UpdateAuctionUseCase defines the interface for updating auctions.
type UpdateAuctionUseCase interface {
	// Execute updates an auction.
	Execute(ctx context.Context, auction *model.Auction) error
}

// updateAuctionUseCase handles updating auctions
type updateAuctionUseCase struct {
	repo repository.AuctionRepository
}

var _ UpdateAuctionUseCase = (*updateAuctionUseCase)(nil)

// NewUpdateAuctionUseCase creates a new instance of UpdateAuctionUseCase
func NewUpdateAuctionUseCase(repo repository.AuctionRepository) *updateAuctionUseCase {
	return &updateAuctionUseCase{repo: repo}
}

// Execute updates an auction
func (uc *updateAuctionUseCase) Execute(ctx context.Context, auction *model.Auction) error {
	return uc.repo.Update(ctx, auction)
}
