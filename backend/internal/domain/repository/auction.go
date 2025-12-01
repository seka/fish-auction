package repository

import (
	"context"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// AuctionFilters represents filters for listing auctions
type AuctionFilters struct {
	VenueID     *int
	AuctionDate *time.Time
	Status      *model.AuctionStatus
	StartDate   *time.Time
	EndDate     *time.Time
}

// AuctionRepository defines the interface for auction data access
type AuctionRepository interface {
	Create(ctx context.Context, auction *model.Auction) (*model.Auction, error)
	GetByID(ctx context.Context, id int) (*model.Auction, error)
	List(ctx context.Context, filters *AuctionFilters) ([]model.Auction, error)
	ListByVenue(ctx context.Context, venueID int) ([]model.Auction, error)
	Update(ctx context.Context, auction *model.Auction) error
	UpdateStatus(ctx context.Context, id int, status model.AuctionStatus) error
	Delete(ctx context.Context, id int) error
}
