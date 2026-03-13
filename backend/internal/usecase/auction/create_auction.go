package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateAuctionUseCase defines the interface for creating auctions
type CreateAuctionUseCase interface {
	Execute(ctx context.Context, auction *model.Auction) (*model.Auction, error)
}

// createAuctionUseCase handles the creation of auctions
type createAuctionUseCase struct {
	repo repository.AuctionRepository
}

var _ CreateAuctionUseCase = (*createAuctionUseCase)(nil)

// NewCreateAuctionUseCase creates a new instance of CreateAuctionUseCase
func NewCreateAuctionUseCase(repo repository.AuctionRepository) *createAuctionUseCase {
	return &createAuctionUseCase{repo: repo}
}

// Execute creates a new auction
func (uc *createAuctionUseCase) Execute(ctx context.Context, auction *model.Auction) (*model.Auction, error) {
	return uc.repo.Create(ctx, auction)
}
