package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockCreateAdminUseCase is a mock implementation of CreateAdminUseCase for testing.
type MockCreateAdminUseCase struct {
	ExecuteFunc func(ctx context.Context, email, password string) (*model.Admin, error)
}

// Execute executes the use case logic.
func (m *MockCreateAdminUseCase) Execute(ctx context.Context, email, password string) (*model.Admin, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, email, password)
	}
	return nil, nil
}

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

// MockAdminVerifyResetTokenUseCase is a mock implementation of VerifyResetTokenUseCase for testing.
type MockAdminVerifyResetTokenUseCase struct {
	ExecuteFunc func(ctx context.Context, token string) error
}

// Execute executes the use case logic.
func (m *MockAdminVerifyResetTokenUseCase) Execute(ctx context.Context, token string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, token)
	}
	return nil
}
