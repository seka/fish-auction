package usecase

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// InvoiceUseCase defines the interface for invoice-related business logic
type InvoiceUseCase interface {
	List(ctx context.Context) ([]model.InvoiceItem, error)
}

type invoiceInteractor struct {
	repo repository.TransactionRepository
}

func NewInvoiceInteractor(repo repository.TransactionRepository) InvoiceUseCase {
	return &invoiceInteractor{repo: repo}
}

func (i *invoiceInteractor) List(ctx context.Context) ([]model.InvoiceItem, error) {
	return i.repo.ListInvoices(ctx)
}
