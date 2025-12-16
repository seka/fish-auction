package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockItemRepository is a mock implementation of ItemRepository
type MockItemRepository struct {
	CreateFunc        func(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
	ListFunc          func(ctx context.Context, status string) ([]model.AuctionItem, error)
	ListByAuctionFunc func(ctx context.Context, auctionID int) ([]model.AuctionItem, error)
	FindByIDFunc      func(ctx context.Context, id int) (*model.AuctionItem, error)
	UpdateStatusFunc  func(ctx context.Context, id int, status model.ItemStatus) error
}

func (m *MockItemRepository) Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	return m.CreateFunc(ctx, item)
}

func (m *MockItemRepository) List(ctx context.Context, status string) ([]model.AuctionItem, error) {
	return m.ListFunc(ctx, status)
}

func (m *MockItemRepository) ListByAuction(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
	return m.ListByAuctionFunc(ctx, auctionID)
}

func (m *MockItemRepository) FindByID(ctx context.Context, id int) (*model.AuctionItem, error) {
	return m.FindByIDFunc(ctx, id)
}

func (m *MockItemRepository) UpdateStatus(ctx context.Context, id int, status model.ItemStatus) error {
	return m.UpdateStatusFunc(ctx, id, status)
}
