package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockBuyerRepository is a mock implementation of BuyerRepository
type MockBuyerRepository struct {
	CreateFunc func(ctx context.Context, name string) (*model.Buyer, error)
	ListFunc   func(ctx context.Context) ([]model.Buyer, error)
}

func (m *MockBuyerRepository) Create(ctx context.Context, name string) (*model.Buyer, error) {
	return m.CreateFunc(ctx, name)
}

func (m *MockBuyerRepository) List(ctx context.Context) ([]model.Buyer, error) {
	return m.ListFunc(ctx)
}
