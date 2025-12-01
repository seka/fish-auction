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
	return uc.repo.GetByID(ctx, id)
}
