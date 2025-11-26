package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockListItemsUseCase is a mock implementation of ListItemsUseCase
type MockListItemsUseCase struct {
	ExecuteFunc func(ctx context.Context, status string) ([]model.AuctionItem, error)
}

func (m *MockListItemsUseCase) Execute(ctx context.Context, status string) ([]model.AuctionItem, error) {
	return m.ExecuteFunc(ctx, status)
}
