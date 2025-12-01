package repository

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// VenueRepository defines the interface for venue data access
type VenueRepository interface {
	Create(ctx context.Context, venue *model.Venue) (*model.Venue, error)
	GetByID(ctx context.Context, id int) (*model.Venue, error)
	List(ctx context.Context) ([]model.Venue, error)
	Update(ctx context.Context, venue *model.Venue) error
	Delete(ctx context.Context, id int) error
}
