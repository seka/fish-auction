package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockFishermanRepository is a mock implementation of FishermanRepository
type MockFishermanRepository struct {
	CreateFunc func(ctx context.Context, name string) (*model.Fisherman, error)
	ListFunc   func(ctx context.Context) ([]model.Fisherman, error)
}

func (m *MockFishermanRepository) Create(ctx context.Context, name string) (*model.Fisherman, error) {
	return m.CreateFunc(ctx, name)
}

func (m *MockFishermanRepository) List(ctx context.Context) ([]model.Fisherman, error) {
	return m.ListFunc(ctx)
}
