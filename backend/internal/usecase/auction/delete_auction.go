package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// DeleteAuctionUseCase defines the interface for deleting auctions
type DeleteAuctionUseCase interface {
	Execute(ctx context.Context, id int) error
}

// deleteAuctionUseCase handles deleting auctions
type deleteAuctionUseCase struct {
	repo repository.AuctionRepository
}

// NewDeleteAuctionUseCase creates a new instance of DeleteAuctionUseCase
func NewDeleteAuctionUseCase(repo repository.AuctionRepository) DeleteAuctionUseCase {
	return &deleteAuctionUseCase{repo: repo}
}

// Execute deletes an auction
func (uc *deleteAuctionUseCase) Execute(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
