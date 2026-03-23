package venue

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// GetVenueUseCase defines the interface for getting a venue by ID
type GetVenueUseCase interface {
	Execute(ctx context.Context, id int) (*model.Venue, error)
}

// GetVenueUseCase handles getting a venue
type getVenueUseCase struct {
	repo repository.VenueRepository
}

var _ GetVenueUseCase = (*getVenueUseCase)(nil)

// NewGetVenueUseCase creates a new instance of GetVenueUseCase
func NewGetVenueUseCase(repo repository.VenueRepository) GetVenueUseCase {
	return &getVenueUseCase{repo: repo}
}

// Execute gets a venue by ID
func (uc *getVenueUseCase) Execute(ctx context.Context, id int) (*model.Venue, error) {
	return uc.repo.FindByID(ctx, id)
}
