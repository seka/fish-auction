package venue

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// UpdateVenueUseCase defines the interface for updating venues
type UpdateVenueUseCase interface {
	Execute(ctx context.Context, venue *model.Venue) error
}

// updateVenueUseCase handles updating venues
type updateVenueUseCase struct {
	repo repository.VenueRepository
}

// NewUpdateVenueUseCase creates a new instance of UpdateVenueUseCase
func NewUpdateVenueUseCase(repo repository.VenueRepository) UpdateVenueUseCase {
	return &updateVenueUseCase{repo: repo}
}

// Execute updates a venue
func (uc *updateVenueUseCase) Execute(ctx context.Context, venue *model.Venue) error {
	return uc.repo.Update(ctx, venue)
}
