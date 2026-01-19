package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockUpdateItemUseCase is a mock implementation of UpdateItemUseCase
type MockUpdateItemUseCase struct {
	ExecuteFunc func(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
}

func (m *MockUpdateItemUseCase) Execute(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, item)
	}
	return item, nil
}

// MockDeleteItemUseCase is a mock implementation of DeleteItemUseCase
type MockDeleteItemUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) error
}

func (m *MockDeleteItemUseCase) Execute(ctx context.Context, id int) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil
}

// MockUpdateItemSortOrderUseCase is a mock implementation of UpdateItemSortOrderUseCase
type MockUpdateItemSortOrderUseCase struct {
	ExecuteFunc func(ctx context.Context, id int, sortOrder int) error
}

func (m *MockUpdateItemSortOrderUseCase) Execute(ctx context.Context, id int, sortOrder int) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id, sortOrder)
	}
	return nil
}
