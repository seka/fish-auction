package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockCreateBidUseCase is a mock implementation of CreateBidUseCase for testing.
type MockCreateBidUseCase struct {
	ExecuteFunc func(ctx context.Context, bid *model.Bid) (*model.Bid, error)
}

// Execute executes the use case logic.
func (m *MockCreateBidUseCase) Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, bid)
	}
	return nil, nil
}
