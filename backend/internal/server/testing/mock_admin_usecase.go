package testing

import (
	"context"
)

type MockAdminUpdatePasswordUseCase struct {
	ExecuteFunc func(ctx context.Context, id int, currentPassword, newPassword string) error
}

func (m *MockAdminUpdatePasswordUseCase) Execute(ctx context.Context, id int, currentPassword, newPassword string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id, currentPassword, newPassword)
	}
	return nil
}

type MockAdminRequestPasswordResetUseCase struct {
	ExecuteFunc func(ctx context.Context, email string) error
}

func (m *MockAdminRequestPasswordResetUseCase) Execute(ctx context.Context, email string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, email)
	}
	return nil
}

type MockAdminResetPasswordUseCase struct {
	ExecuteFunc func(ctx context.Context, token, newPassword string) error
}

func (m *MockAdminResetPasswordUseCase) Execute(ctx context.Context, token, newPassword string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, token, newPassword)
	}
	return nil
}
