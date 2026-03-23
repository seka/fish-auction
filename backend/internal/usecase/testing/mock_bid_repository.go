package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockBidRepository is a mock implementation of BidRepository
type MockBidRepository struct {
	CreateFunc                 func(ctx context.Context, bid *model.Bid) (*model.Bid, error)
	ListInvoicesFunc           func(ctx context.Context) ([]model.InvoiceItem, error)
	ListPurchasesByBuyerIDFunc func(ctx context.Context, buyerID int) ([]model.Purchase, error)
	ListAuctionsByBuyerIDFunc  func(ctx context.Context, buyerID int) ([]model.Auction, error)
}

// Create creates a new record.
func (m *MockBidRepository) Create(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	return m.CreateFunc(ctx, bid)
}

// ListInvoices retrieves a list of records.
func (m *MockBidRepository) ListInvoices(ctx context.Context) ([]model.InvoiceItem, error) {
	return m.ListInvoicesFunc(ctx)
}

// ListPurchasesByBuyerID retrieves a list of records.
func (m *MockBidRepository) ListPurchasesByBuyerID(ctx context.Context, buyerID int) ([]model.Purchase, error) {
	return m.ListPurchasesByBuyerIDFunc(ctx, buyerID)
}

// ListAuctionsByBuyerID retrieves a list of records.
func (m *MockBidRepository) ListAuctionsByBuyerID(ctx context.Context, buyerID int) ([]model.Auction, error) {
	return m.ListAuctionsByBuyerIDFunc(ctx, buyerID)
}
