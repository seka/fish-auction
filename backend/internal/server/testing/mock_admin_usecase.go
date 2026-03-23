package testing

import (
	"context"
)

// MockAdminUpdatePasswordUseCase is a mock implementation of AdminUpdatePasswordUseCase for testing.
type MockAdminUpdatePasswordUseCase struct {
	ExecuteFunc func(ctx context.Context, id int, currentPassword, newPassword string) error
}

// Execute executes the use case logic.
func (m *MockAdminUpdatePasswordUseCase) Execute(ctx context.Context, id int, currentPassword, newPassword string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, id, currentPassword, newPassword)
	}
	return nil
}

// MockAdminRequestPasswordResetUseCase is a mock implementation of AdminRequestPasswordResetUseCase for testing.
type MockAdminRequestPasswordResetUseCase struct {
	ExecuteFunc func(ctx context.Context, email string) error
}

// Execute executes the use case logic.
func (m *MockAdminRequestPasswordResetUseCase) Execute(ctx context.Context, email string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, email)
	}
	return nil
}

// MockAdminResetPasswordUseCase is a mock implementation of AdminResetPasswordUseCase for testing.
type MockAdminResetPasswordUseCase struct {
	ExecuteFunc func(ctx context.Context, token, newPassword string) error
}

// Execute executes the use case logic.
func (m *MockAdminResetPasswordUseCase) Execute(ctx context.Context, token, newPassword string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, token, newPassword)
	}
	return nil
}
