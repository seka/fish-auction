package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockBuyerRepository is a mock implementation of BuyerRepository
type MockBuyerRepository struct {
	CreateFunc     func(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error)
	ListFunc       func(ctx context.Context) ([]model.Buyer, error)
	FindByIDFunc   func(ctx context.Context, id int) (*model.Buyer, error)
	FindByNameFunc func(ctx context.Context, name string) (*model.Buyer, error)
}

func (m *MockBuyerRepository) Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error) {
	return m.CreateFunc(ctx, buyer)
}

func (m *MockBuyerRepository) List(ctx context.Context) ([]model.Buyer, error) {
	return m.ListFunc(ctx)
}

func (m *MockBuyerRepository) FindByID(ctx context.Context, id int) (*model.Buyer, error) {
	return m.FindByIDFunc(ctx, id)
}

func (m *MockBuyerRepository) FindByName(ctx context.Context, name string) (*model.Buyer, error) {
	return m.FindByNameFunc(ctx, name)
}
