package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockCreateFishermanUseCase is a mock implementation of CreateFishermanUseCase for testing.
type MockCreateFishermanUseCase struct {
	ExecuteFunc func(ctx context.Context, name string) (*model.Fisherman, error)
}

// Execute executes the use case logic.
func (m *MockCreateFishermanUseCase) Execute(ctx context.Context, name string) (*model.Fisherman, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, name)
	}
	return nil, nil
}

// MockListFishermenUseCase is a mock implementation of ListFishermenUseCase for testing.
type MockListFishermenUseCase struct {
	ExecuteFunc func(ctx context.Context) ([]model.Fisherman, error)
}

// Execute executes the use case logic.
func (m *MockListFishermenUseCase) Execute(ctx context.Context) ([]model.Fisherman, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx)
	}
	return nil, nil
}

// MockDeleteFishermanUseCase is a mock implementation of DeleteFishermanUseCase for testing.
type MockDeleteFishermanUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) error
}

// Execute executes the use case logic.
func (m *MockDeleteFishermanUseCase) Execute(ctx context.Context, id int) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil
}
