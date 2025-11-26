package invoice

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListInvoicesUseCase handles listing invoices
type ListInvoicesUseCase struct {
	bidRepo repository.BidRepository
}

// NewListInvoicesUseCase creates a new instance of ListInvoicesUseCase
func NewListInvoicesUseCase(bidRepo repository.BidRepository) *ListInvoicesUseCase {
	return &ListInvoicesUseCase{bidRepo: bidRepo}
}

// Execute lists all invoices grouped by buyer
func (uc *ListInvoicesUseCase) Execute(ctx context.Context) ([]model.InvoiceItem, error) {
	return uc.bidRepo.ListInvoices(ctx)
}
