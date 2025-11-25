package repository

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type TransactionRepository interface {
	Create(ctx context.Context, t *model.Transaction) (*model.Transaction, error)
	ListInvoices(ctx context.Context) ([]model.InvoiceItem, error)
}
