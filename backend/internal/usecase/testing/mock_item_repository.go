package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockItemRepository is a mock implementation of ItemRepository
type MockItemRepository struct {
	CreateFunc           func(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
	ListFunc             func(ctx context.Context) ([]model.AuctionItem, error)
	ListByAuctionFunc    func(ctx context.Context, auctionID int) ([]model.AuctionItem, error)
	FindByIDFunc         func(ctx context.Context, id int) (*model.AuctionItem, error)
	FindByIDWithLockFunc func(ctx context.Context, id int) (*model.AuctionItem, error)
	UpdateFunc           func(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
	DeleteFunc           func(ctx context.Context, id int) error
	UpdateSortOrderFunc  func(ctx context.Context, id int, sortOrder int) error
	ReorderFunc          func(ctx context.Context, auctionID int, ids []int) error
}

// Create creates a new record.
func (m *MockItemRepository) Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	return m.CreateFunc(ctx, item)
}

// List retrieves a list of records.
func (m *MockItemRepository) List(ctx context.Context) ([]model.AuctionItem, error) {
	return m.ListFunc(ctx)
}

// ListByAuction retrieves a list of records.
func (m *MockItemRepository) ListByAuction(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
	return m.ListByAuctionFunc(ctx, auctionID)
}

// FindByID retrieves a record based on criteria.
func (m *MockItemRepository) FindByID(ctx context.Context, id int) (*model.AuctionItem, error) {
	return m.FindByIDFunc(ctx, id)
}

// FindByIDWithLock retrieves a record based on criteria.
func (m *MockItemRepository) FindByIDWithLock(ctx context.Context, id int) (*model.AuctionItem, error) {
	return m.FindByIDWithLockFunc(ctx, id)
}

// Update updates an existing record.
func (m *MockItemRepository) Update(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	return m.UpdateFunc(ctx, item)
}

// Delete removes a record by ID.
func (m *MockItemRepository) Delete(ctx context.Context, id int) error {
	return m.DeleteFunc(ctx, id)
}

// UpdateSortOrder updates an existing record.
func (m *MockItemRepository) UpdateSortOrder(ctx context.Context, id, sortOrder int) error {
	return m.UpdateSortOrderFunc(ctx, id, sortOrder)
}

// Reorder provides Reorder related functionality.
func (m *MockItemRepository) Reorder(ctx context.Context, auctionID int, ids []int) error {
	return m.ReorderFunc(ctx, auctionID, ids)
}
