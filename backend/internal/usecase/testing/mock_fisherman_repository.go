package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockFishermanRepository is a mock implementation of FishermanRepository
type MockFishermanRepository struct {
	CreateFunc   func(ctx context.Context, name string) (*model.Fisherman, error)
	ListFunc     func(ctx context.Context) ([]model.Fisherman, error)
	FindByIDFunc func(ctx context.Context, id int) (*model.Fisherman, error)
	DeleteFunc   func(ctx context.Context, id int) error
}

// Create creates a new record.
func (m *MockFishermanRepository) Create(ctx context.Context, name string) (*model.Fisherman, error) {
	return m.CreateFunc(ctx, name)
}

// List retrieves a list of records.
func (m *MockFishermanRepository) List(ctx context.Context) ([]model.Fisherman, error) {
	return m.ListFunc(ctx)
}

// FindByID retrieves a record based on criteria.
func (m *MockFishermanRepository) FindByID(ctx context.Context, id int) (*model.Fisherman, error) {
	return m.FindByIDFunc(ctx, id)
}

// Delete removes a record by ID.
func (m *MockFishermanRepository) Delete(ctx context.Context, id int) error {
	return m.DeleteFunc(ctx, id)
}
