package venue

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListVenuesUseCase defines the interface for listing venues
type ListVenuesUseCase interface {
	Execute(ctx context.Context) ([]model.Venue, error)
}

// listVenuesUseCase handles listing venues
type listVenuesUseCase struct {
	repo repository.VenueRepository
}

// NewListVenuesUseCase creates a new instance of ListVenuesUseCase
func NewListVenuesUseCase(repo repository.VenueRepository) ListVenuesUseCase {
	return &listVenuesUseCase{repo: repo}
}

// Execute lists all venues
func (uc *listVenuesUseCase) Execute(ctx context.Context) ([]model.Venue, error) {
	return uc.repo.List(ctx)
}
