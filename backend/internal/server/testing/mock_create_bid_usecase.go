package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockCreateBidUseCase is a mock implementation of CreateBidUseCase
type MockCreateBidUseCase struct {
	ExecuteFunc func(ctx context.Context, bid *model.Bid) (*model.Bid, error)
}

func (m *MockCreateBidUseCase) Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	return m.ExecuteFunc(ctx, bid)
}
