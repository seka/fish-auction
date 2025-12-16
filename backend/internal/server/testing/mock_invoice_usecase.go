package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type MockListInvoicesUseCase struct {
	ExecuteFunc func(ctx context.Context) ([]model.InvoiceItem, error)
}

func (m *MockListInvoicesUseCase) Execute(ctx context.Context) ([]model.InvoiceItem, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx)
	}
	return nil, nil
}
