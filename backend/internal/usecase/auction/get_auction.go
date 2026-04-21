package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// GetAuctionUseCase defines the interface for getting an auction by ID.
type GetAuctionUseCase interface {
	// Execute gets an auction by ID.
	Execute(ctx context.Context, id int) (*model.Auction, error)
}

// GetAuctionUseCase handles getting an auction
type getAuctionUseCase struct {
	repo repository.AuctionRepository
}

var _ GetAuctionUseCase = (*getAuctionUseCase)(nil)

// NewGetAuctionUseCase creates a new instance of GetAuctionUseCase
func NewGetAuctionUseCase(auctionRepo repository.AuctionRepository) GetAuctionUseCase {
	return &getAuctionUseCase{
		repo: auctionRepo,
	}
}

// Execute gets an auction by ID
func (uc *getAuctionUseCase) Execute(ctx context.Context, id int) (*model.Auction, error) {
	return uc.repo.FindByID(ctx, id)
}
