package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockCreateVenueUseCase is a mock implementation of CreateVenueUseCase for testing.
type MockCreateVenueUseCase struct {
	ExecuteFunc func(ctx context.Context, venue *model.Venue) (*model.Venue, error)
}

// Execute executes the use case logic.
func (m *MockCreateVenueUseCase) Execute(ctx context.Context, venue *model.Venue) (*model.Venue, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, venue)
	}
	return nil, nil
}

// MockListVenuesUseCase is a mock implementation of ListVenuesUseCase for testing.
type MockListVenuesUseCase struct {
	ExecuteFunc func(ctx context.Context) ([]model.Venue, error)
}

// Execute executes the use case logic.
func (m *MockListVenuesUseCase) Execute(ctx context.Context) ([]model.Venue, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx)
	}
	return nil, nil
}

// MockGetVenueUseCase is a mock implementation of GetVenueUseCase for testing.
type MockGetVenueUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) (*model.Venue, error)
}

// Execute executes the use case logic.
func (m *MockGetVenueUseCase) Execute(ctx context.Context, id int) (*model.Venue, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil, nil
}

// MockUpdateVenueUseCase is a mock implementation of UpdateVenueUseCase for testing.
type MockUpdateVenueUseCase struct {
	ExecuteFunc func(ctx context.Context, venue *model.Venue) error
}

// Execute executes the use case logic.
func (m *MockUpdateVenueUseCase) Execute(ctx context.Context, venue *model.Venue) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, venue)
	}
	return nil
}

// MockDeleteVenueUseCase is a mock implementation of DeleteVenueUseCase for testing.
type MockDeleteVenueUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) error
}

// Execute executes the use case logic.
func (m *MockDeleteVenueUseCase) Execute(ctx context.Context, id int) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil
}
