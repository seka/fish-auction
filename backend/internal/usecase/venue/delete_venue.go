package venue

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// DeleteVenueUseCase defines the interface for deleting venues
type DeleteVenueUseCase interface {
	Execute(ctx context.Context, id int) error
}

// deleteVenueUseCase handles deleting venues
type deleteVenueUseCase struct {
	repo repository.VenueRepository
}

var _ DeleteVenueUseCase = (*deleteVenueUseCase)(nil)

// NewDeleteVenueUseCase creates a new instance of DeleteVenueUseCase
func NewDeleteVenueUseCase(repo repository.VenueRepository) *deleteVenueUseCase {
	return &deleteVenueUseCase{repo: repo}
}

// Execute deletes a venue
func (uc *deleteVenueUseCase) Execute(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
