package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
)

type MockLoginUseCase struct {
	ExecuteFunc func(ctx context.Context, email, password string) (*entity.Admin, error)
}

func (m *MockLoginUseCase) Execute(ctx context.Context, email, password string) (*entity.Admin, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, email, password)
	}
	return nil, nil
}

type MockRequestPasswordResetUseCase struct {
	ExecuteFunc func(ctx context.Context, email string) error
}

func (m *MockRequestPasswordResetUseCase) Execute(ctx context.Context, email string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, email)
	}
	return nil
}

type MockResetPasswordUseCase struct {
	ExecuteFunc func(ctx context.Context, token, newPassword string) error
}

func (m *MockResetPasswordUseCase) Execute(ctx context.Context, token, newPassword string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, token, newPassword)
	}
	return nil
}
