package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type MockCreateAuctionUseCase struct {
	ExecuteFunc func(ctx context.Context, auction *model.Auction) (*model.Auction, error)
}

func (m *MockCreateAuctionUseCase) Execute(ctx context.Context, auction *model.Auction) (*model.Auction, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, auction)
	}
	return nil, nil
}

type MockListAuctionsUseCase struct {
	ExecuteFunc func(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error)
}

func (m *MockListAuctionsUseCase) Execute(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, filters)
	}
	return nil, nil
}

type MockGetAuctionUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) (*model.Auction, error)
}

func (m *MockGetAuctionUseCase) Execute(ctx context.Context, id int) (*model.Auction, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil, nil
}

type MockGetAuctionItemsUseCase struct {
	ExecuteFunc func(ctx context.Context, auctionID int) ([]model.AuctionItem, error)
}

func (m *MockGetAuctionItemsUseCase) Execute(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, auctionID)
	}
	return nil, nil
}

type MockUpdateAuctionUseCase struct {
	ExecuteFunc func(ctx context.Context, auction *model.Auction) error
}

func (m *MockUpdateAuctionUseCase) Execute(ctx context.Context, auction *model.Auction) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, auction)
	}
	return nil
}

type MockUpdateAuctionStatusUseCase struct {
	ExecuteFunc func(ctx context.Context, id int, status model.AuctionStatus) error
}

func (m *MockUpdateAuctionStatusUseCase) Execute(ctx context.Context, id int, status model.AuctionStatus) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id, status)
	}
	return nil
}

type MockDeleteAuctionUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) error
}

func (m *MockDeleteAuctionUseCase) Execute(ctx context.Context, id int) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil
}
