package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockCreateBuyerUseCase is a mock implementation of CreateBuyerUseCase for testing.
type MockCreateBuyerUseCase struct {
	ExecuteFunc func(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error)
}

// Execute executes the use case logic.
func (m *MockCreateBuyerUseCase) Execute(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, name, email, password, organization, contactInfo)
	}
	return nil, nil
}

// MockListBuyersUseCase is a mock implementation of ListBuyersUseCase for testing.
type MockListBuyersUseCase struct {
	ExecuteFunc func(ctx context.Context) ([]model.Buyer, error)
}

// Execute executes the use case logic.
func (m *MockListBuyersUseCase) Execute(ctx context.Context) ([]model.Buyer, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx)
	}
	return nil, nil
}

// MockLoginBuyerUseCase is a mock implementation of LoginBuyerUseCase for testing.
type MockLoginBuyerUseCase struct {
	ExecuteFunc func(ctx context.Context, email, password string) (*model.Buyer, error)
}

// Execute executes the use case logic.
func (m *MockLoginBuyerUseCase) Execute(ctx context.Context, email, password string) (*model.Buyer, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, email, password)
	}
	return nil, nil
}

// MockGetBuyerPurchasesUseCase is a mock implementation of GetBuyerPurchasesUseCase for testing.
type MockGetBuyerPurchasesUseCase struct {
	ExecuteFunc func(ctx context.Context, buyerID int) ([]model.Purchase, error)
}

// Execute executes the use case logic.
func (m *MockGetBuyerPurchasesUseCase) Execute(ctx context.Context, buyerID int) ([]model.Purchase, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, buyerID)
	}
	return nil, nil
}

// MockGetBuyerAuctionsUseCase is a mock implementation of GetBuyerAuctionsUseCase for testing.
type MockGetBuyerAuctionsUseCase struct {
	ExecuteFunc func(ctx context.Context, buyerID int) ([]model.Auction, error)
}

// Execute executes the use case logic.
func (m *MockGetBuyerAuctionsUseCase) Execute(ctx context.Context, buyerID int) ([]model.Auction, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, buyerID)
	}
	return nil, nil
}

// MockBuyerUpdatePasswordUseCase is a mock implementation of BuyerUpdatePasswordUseCase for testing.
type MockBuyerUpdatePasswordUseCase struct {
	ExecuteFunc func(ctx context.Context, buyerID int, currentPassword, newPassword string) error
}

// Execute executes the use case logic.
func (m *MockBuyerUpdatePasswordUseCase) Execute(ctx context.Context, buyerID int, currentPassword, newPassword string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, buyerID, currentPassword, newPassword)
	}
	return nil
}

// MockGetBuyerUseCase is a mock implementation of GetBuyerUseCase for testing.
type MockGetBuyerUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) (*model.Buyer, error)
}

// Execute executes the use case logic.
func (m *MockGetBuyerUseCase) Execute(ctx context.Context, id int) (*model.Buyer, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil, nil
}

// MockDeleteBuyerUseCase is a mock implementation of DeleteBuyerUseCase for testing.
type MockDeleteBuyerUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) error
}

// Execute executes the use case logic.
func (m *MockDeleteBuyerUseCase) Execute(ctx context.Context, id int) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil
}
