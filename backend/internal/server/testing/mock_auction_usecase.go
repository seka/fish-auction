package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// MockCreateAuctionUseCase is a mock implementation of CreateAuctionUseCase for testing.
type MockCreateAuctionUseCase struct {
	ExecuteFunc func(ctx context.Context, auction *model.Auction) (*model.Auction, error)
}

// Execute executes the use case logic.
func (m *MockCreateAuctionUseCase) Execute(ctx context.Context, auction *model.Auction) (*model.Auction, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, auction)
	}
	return nil, nil
}

// MockListAuctionsUseCase is a mock implementation of ListAuctionsUseCase for testing.
type MockListAuctionsUseCase struct {
	ExecuteFunc func(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error)
}

// Execute executes the use case logic.
func (m *MockListAuctionsUseCase) Execute(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, filters)
	}
	return nil, nil
}

// MockGetAuctionUseCase is a mock implementation of GetAuctionUseCase for testing.
type MockGetAuctionUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) (*model.Auction, error)
}

// Execute executes the use case logic.
func (m *MockGetAuctionUseCase) Execute(ctx context.Context, id int) (*model.Auction, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil, nil
}

// MockGetAuctionItemsUseCase is a mock implementation of GetAuctionItemsUseCase for testing.
type MockGetAuctionItemsUseCase struct {
	ExecuteFunc func(ctx context.Context, auctionID int) ([]model.AuctionItem, error)
}

// Execute executes the use case logic.
func (m *MockGetAuctionItemsUseCase) Execute(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, auctionID)
	}
	return nil, nil
}

// MockUpdateAuctionUseCase is a mock implementation of UpdateAuctionUseCase for testing.
type MockUpdateAuctionUseCase struct {
	ExecuteFunc func(ctx context.Context, auction *model.Auction) error
}

// Execute executes the use case logic.
func (m *MockUpdateAuctionUseCase) Execute(ctx context.Context, auction *model.Auction) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, auction)
	}
	return nil
}

// MockUpdateAuctionStatusUseCase is a mock implementation of UpdateAuctionStatusUseCase for testing.
type MockUpdateAuctionStatusUseCase struct {
	ExecuteFunc func(ctx context.Context, id int, status model.AuctionStatus) error
}

// Execute executes the use case logic.
func (m *MockUpdateAuctionStatusUseCase) Execute(ctx context.Context, id int, status model.AuctionStatus) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id, status)
	}
	return nil
}

// MockDeleteAuctionUseCase is a mock implementation of DeleteAuctionUseCase for testing.
type MockDeleteAuctionUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) error
}

// Execute executes the use case logic.
func (m *MockDeleteAuctionUseCase) Execute(ctx context.Context, id int) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil
}
