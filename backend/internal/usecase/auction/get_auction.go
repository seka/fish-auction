package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// GetAuctionUseCase defines the interface for getting an auction by ID
type GetAuctionUseCase interface {
	Execute(ctx context.Context, id int) (*model.Auction, error)
}

// getAuctionUseCase handles getting an auction
type getAuctionUseCase struct {
	repo repository.AuctionRepository
}

// NewGetAuctionUseCase creates a new instance of GetAuctionUseCase
func NewGetAuctionUseCase(repo repository.AuctionRepository) GetAuctionUseCase {
	return &getAuctionUseCase{repo: repo}
}

// Execute gets an auction by ID
func (uc *getAuctionUseCase) Execute(ctx context.Context, id int) (*model.Auction, error) {
	auction, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if auction.ShouldBeCompleted() {
		if err := uc.repo.UpdateStatus(ctx, auction.ID, model.AuctionStatusCompleted); err != nil {
			return nil, err
		}
		auction.Status = model.AuctionStatusCompleted
	}

	return auction, nil
}
