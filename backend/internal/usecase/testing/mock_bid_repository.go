package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockBidRepository is a mock implementation of BidRepository
type MockBidRepository struct {
	CreateFunc       func(ctx context.Context, bid *model.Bid) (*model.Bid, error)
	ListInvoicesFunc func(ctx context.Context) ([]model.InvoiceItem, error)
}

func (m *MockBidRepository) Create(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	return m.CreateFunc(ctx, bid)
}

func (m *MockBidRepository) ListInvoices(ctx context.Context) ([]model.InvoiceItem, error) {
	return m.ListInvoicesFunc(ctx)
}
