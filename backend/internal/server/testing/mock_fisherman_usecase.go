package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type MockCreateFishermanUseCase struct {
	ExecuteFunc func(ctx context.Context, name string) (*model.Fisherman, error)
}

func (m *MockCreateFishermanUseCase) Execute(ctx context.Context, name string) (*model.Fisherman, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, name)
	}
	return nil, nil
}

type MockListFishermenUseCase struct {
	ExecuteFunc func(ctx context.Context) ([]model.Fisherman, error)
}

func (m *MockListFishermenUseCase) Execute(ctx context.Context) ([]model.Fisherman, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx)
	}
	return nil, nil
}
