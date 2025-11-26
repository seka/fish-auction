package invoice

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListInvoicesUseCase defines the interface for listing invoices
type ListInvoicesUseCase interface {
	Execute(ctx context.Context) ([]model.InvoiceItem, error)
}

// listInvoicesUseCase handles listing invoices
type listInvoicesUseCase struct {
	bidRepo repository.BidRepository
}

// NewListInvoicesUseCase creates a new instance of ListInvoicesUseCase
func NewListInvoicesUseCase(bidRepo repository.BidRepository) ListInvoicesUseCase {
	return &listInvoicesUseCase{bidRepo: bidRepo}
}

// Execute lists all invoices grouped by buyer
func (uc *listInvoicesUseCase) Execute(ctx context.Context) ([]model.InvoiceItem, error) {
	return uc.bidRepo.ListInvoices(ctx)
}
