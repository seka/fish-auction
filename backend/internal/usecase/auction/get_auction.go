package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

// GetAuctionUseCase defines the interface for getting an auction by ID.
type GetAuctionUseCase interface {
	// Execute gets an auction by ID.
	Execute(ctx context.Context, id int) (*model.Auction, error)
}

// GetAuctionUseCase handles getting an auction
type getAuctionUseCase struct {
	repo  repository.AuctionRepository
	clock service.Clock
}

var _ GetAuctionUseCase = (*getAuctionUseCase)(nil)

// NewGetAuctionUseCase creates a new instance of GetAuctionUseCase
func NewGetAuctionUseCase(auctionRepo repository.AuctionRepository, clock service.Clock) GetAuctionUseCase {
	return &getAuctionUseCase{
		repo:  auctionRepo,
		clock: clock,
	}
}

// Execute gets an auction by ID
func (uc *getAuctionUseCase) Execute(ctx context.Context, id int) (*model.Auction, error) {
	auction, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	now := uc.clock.Now()
	if auction.ShouldBeCompleted(now) {
		if err := uc.repo.UpdateStatus(ctx, auction.ID, model.AuctionStatusCompleted); err != nil {
			return nil, err
		}
		auction.Status = model.AuctionStatusCompleted
	}

	return auction, nil
}
