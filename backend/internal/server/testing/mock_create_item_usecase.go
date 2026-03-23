package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockCreateItemUseCase is a mock implementation of CreateItemUseCase
type MockCreateItemUseCase struct {
	ExecuteFunc func(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
}

// Execute executes the use case logic.
func (m *MockCreateItemUseCase) Execute(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	return m.ExecuteFunc(ctx, item)
}
