package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type MockCreateBuyerUseCase struct {
	ExecuteFunc func(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error)
}

func (m *MockCreateBuyerUseCase) Execute(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, name, email, password, organization, contactInfo)
	}
	return nil, nil
}

type MockListBuyersUseCase struct {
	ExecuteFunc func(ctx context.Context) ([]model.Buyer, error)
}

func (m *MockListBuyersUseCase) Execute(ctx context.Context) ([]model.Buyer, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx)
	}
	return nil, nil
}

type MockLoginBuyerUseCase struct {
	ExecuteFunc func(ctx context.Context, email, password string) (*model.Buyer, error)
}

func (m *MockLoginBuyerUseCase) Execute(ctx context.Context, email, password string) (*model.Buyer, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, email, password)
	}
	return nil, nil
}

type MockGetBuyerPurchasesUseCase struct {
	ExecuteFunc func(ctx context.Context, buyerID int) ([]model.Purchase, error)
}

func (m *MockGetBuyerPurchasesUseCase) Execute(ctx context.Context, buyerID int) ([]model.Purchase, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, buyerID)
	}
	return nil, nil
}

type MockGetBuyerAuctionsUseCase struct {
	ExecuteFunc func(ctx context.Context, buyerID int) ([]model.Auction, error)
}

func (m *MockGetBuyerAuctionsUseCase) Execute(ctx context.Context, buyerID int) ([]model.Auction, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, buyerID)
	}
	return nil, nil
}

type MockBuyerUpdatePasswordUseCase struct {
	ExecuteFunc func(ctx context.Context, buyerID int, currentPassword, newPassword string) error
}

func (m *MockBuyerUpdatePasswordUseCase) Execute(ctx context.Context, buyerID int, currentPassword, newPassword string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, buyerID, currentPassword, newPassword)
	}
	return nil
}
