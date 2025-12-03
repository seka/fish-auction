package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// MockAuctionRepository is a mock implementation of repository.AuctionRepository
type MockAuctionRepository struct {
	CreateFunc       func(ctx context.Context, auction *model.Auction) (*model.Auction, error)
	GetByIDFunc      func(ctx context.Context, id int) (*model.Auction, error)
	ListFunc         func(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error)
	ListByVenueFunc  func(ctx context.Context, venueID int) ([]model.Auction, error)
	UpdateFunc       func(ctx context.Context, auction *model.Auction) error
	UpdateStatusFunc func(ctx context.Context, id int, status model.AuctionStatus) error
	DeleteFunc       func(ctx context.Context, id int) error
}

func (m *MockAuctionRepository) Create(ctx context.Context, auction *model.Auction) (*model.Auction, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, auction)
	}
	return nil, nil
}

func (m *MockAuctionRepository) GetByID(ctx context.Context, id int) (*model.Auction, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockAuctionRepository) List(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, filters)
	}
	return nil, nil
}

func (m *MockAuctionRepository) ListByVenue(ctx context.Context, venueID int) ([]model.Auction, error) {
	if m.ListByVenueFunc != nil {
		return m.ListByVenueFunc(ctx, venueID)
	}
	return nil, nil
}

func (m *MockAuctionRepository) Update(ctx context.Context, auction *model.Auction) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, auction)
	}
	return nil
}

func (m *MockAuctionRepository) UpdateStatus(ctx context.Context, id int, status model.AuctionStatus) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *MockAuctionRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
