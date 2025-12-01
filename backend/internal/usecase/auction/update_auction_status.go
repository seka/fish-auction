package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// UpdateAuctionStatusUseCase defines the interface for updating auction status
type UpdateAuctionStatusUseCase interface {
	Execute(ctx context.Context, id int, status model.AuctionStatus) error
}

// updateAuctionStatusUseCase handles updating auction status
type updateAuctionStatusUseCase struct {
	repo repository.AuctionRepository
}

// NewUpdateAuctionStatusUseCase creates a new instance of UpdateAuctionStatusUseCase
func NewUpdateAuctionStatusUseCase(repo repository.AuctionRepository) UpdateAuctionStatusUseCase {
	return &updateAuctionStatusUseCase{repo: repo}
}

// Execute updates an auction's status
func (uc *updateAuctionStatusUseCase) Execute(ctx context.Context, id int, status model.AuctionStatus) error {
	// Validate status
	if !status.IsValid() {
		return &InvalidStatusError{Status: string(status)}
	}
	return uc.repo.UpdateStatus(ctx, id, status)
}

// InvalidStatusError represents an invalid auction status error
type InvalidStatusError struct {
	Status string
}

func (e *InvalidStatusError) Error() string {
	return "invalid auction status: " + e.Status
}
