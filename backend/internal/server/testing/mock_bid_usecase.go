package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type MockCreateBidUseCase struct {
	ExecuteFunc func(ctx context.Context, bid *model.Bid) (*model.Bid, error)
}

func (m *MockCreateBidUseCase) Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, bid)
	}
	return nil, nil
}
