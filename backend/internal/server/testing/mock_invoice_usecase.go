package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockListInvoicesUseCase is a mock implementation of ListInvoicesUseCase for testing.
type MockListInvoicesUseCase struct {
	ExecuteFunc func(ctx context.Context) ([]model.InvoiceItem, error)
}

// Execute executes the use case logic.
func (m *MockListInvoicesUseCase) Execute(ctx context.Context) ([]model.InvoiceItem, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx)
	}
	return nil, nil
}
