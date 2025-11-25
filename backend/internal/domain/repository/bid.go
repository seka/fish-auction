package repository

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type BidRepository interface {
	Create(ctx context.Context, bid *model.Bid) (*model.Bid, error)
	ListInvoices(ctx context.Context) ([]model.InvoiceItem, error)
}
