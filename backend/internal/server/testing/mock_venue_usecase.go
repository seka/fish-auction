package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type MockCreateVenueUseCase struct {
	ExecuteFunc func(ctx context.Context, venue *model.Venue) (*model.Venue, error)
}

func (m *MockCreateVenueUseCase) Execute(ctx context.Context, venue *model.Venue) (*model.Venue, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, venue)
	}
	return nil, nil
}

type MockListVenuesUseCase struct {
	ExecuteFunc func(ctx context.Context) ([]model.Venue, error)
}

func (m *MockListVenuesUseCase) Execute(ctx context.Context) ([]model.Venue, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx)
	}
	return nil, nil
}

type MockGetVenueUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) (*model.Venue, error)
}

func (m *MockGetVenueUseCase) Execute(ctx context.Context, id int) (*model.Venue, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil, nil
}

type MockUpdateVenueUseCase struct {
	ExecuteFunc func(ctx context.Context, venue *model.Venue) error
}

func (m *MockUpdateVenueUseCase) Execute(ctx context.Context, venue *model.Venue) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, venue)
	}
	return nil
}

type MockDeleteVenueUseCase struct {
	ExecuteFunc func(ctx context.Context, id int) error
}

func (m *MockDeleteVenueUseCase) Execute(ctx context.Context, id int) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id)
	}
	return nil
}
