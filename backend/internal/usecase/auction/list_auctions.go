package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

// ListAuctionsUseCase defines the interface for listing auctions.
type ListAuctionsUseCase interface {
	// Execute lists auctions with optional filters.
	Execute(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error)
}

// ListAuctionsUseCase handles listing auctions
type listAuctionsUseCase struct {
	repo  repository.AuctionRepository
	clock service.Clock
}

var _ ListAuctionsUseCase = (*listAuctionsUseCase)(nil)

// NewListAuctionsUseCase creates a new instance of ListAuctionsUseCase
func NewListAuctionsUseCase(auctionRepo repository.AuctionRepository, clock service.Clock) ListAuctionsUseCase {
	return &listAuctionsUseCase{
		repo:  auctionRepo,
		clock: clock,
	}
}

// Execute lists auctions with optional filters
func (uc *listAuctionsUseCase) Execute(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
	auctions, err := uc.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	now := uc.clock.Now()
	for i := range auctions {
		if auctions[i].ShouldBeCompleted(now) {
			// Update status in-memory only for display purposes
			// This prevents N+1 UPDATE queries on every list request
			auctions[i].Status = model.AuctionStatusCompleted
		}
	}

	return auctions, nil
}
