package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListAuctionsUseCase defines the interface for listing auctions
type ListAuctionsUseCase interface {
	Execute(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error)
}

// listAuctionsUseCase handles listing auctions
type listAuctionsUseCase struct {
	repo repository.AuctionRepository
}

// NewListAuctionsUseCase creates a new instance of ListAuctionsUseCase
func NewListAuctionsUseCase(repo repository.AuctionRepository) ListAuctionsUseCase {
	return &listAuctionsUseCase{repo: repo}
}

// Execute lists auctions with optional filters
func (uc *listAuctionsUseCase) Execute(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
	auctions, err := uc.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	for i := range auctions {
		if auctions[i].ShouldBeCompleted() {
			if err := uc.repo.UpdateStatus(ctx, auctions[i].ID, model.AuctionStatusCompleted); err == nil {
				auctions[i].Status = model.AuctionStatusCompleted
			}
		}
	}

	return auctions, nil
}
