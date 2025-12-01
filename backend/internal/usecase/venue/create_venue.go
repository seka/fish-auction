package venue

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateVenueUseCase defines the interface for creating venues
type CreateVenueUseCase interface {
	Execute(ctx context.Context, venue *model.Venue) (*model.Venue, error)
}

// createVenueUseCase handles the creation of venues
type createVenueUseCase struct {
	repo repository.VenueRepository
}

// NewCreateVenueUseCase creates a new instance of CreateVenueUseCase
func NewCreateVenueUseCase(repo repository.VenueRepository) CreateVenueUseCase {
	return &createVenueUseCase{repo: repo}
}

// Execute creates a new venue
func (uc *createVenueUseCase) Execute(ctx context.Context, venue *model.Venue) (*model.Venue, error) {
	return uc.repo.Create(ctx, venue)
}
