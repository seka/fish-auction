package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockLoginUseCase is a mock implementation of LoginUseCase for testing.
type MockLoginUseCase struct {
	ExecuteFunc func(ctx context.Context, email, password string) (*model.Admin, error)
}

// Execute executes the use case logic.
func (m *MockLoginUseCase) Execute(ctx context.Context, email, password string) (*model.Admin, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, email, password)
	}
	return nil, nil
}

// MockRequestPasswordResetUseCase is a mock implementation of RequestPasswordResetUseCase for testing.
type MockRequestPasswordResetUseCase struct {
	ExecuteFunc func(ctx context.Context, email string) error
}

// Execute executes the use case logic.
func (m *MockRequestPasswordResetUseCase) Execute(ctx context.Context, email string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, email)
	}
	return nil
}

// MockResetPasswordUseCase is a mock implementation of ResetPasswordUseCase for testing.
type MockResetPasswordUseCase struct {
	ExecuteFunc func(ctx context.Context, token, newPassword string) error
}

// Execute executes the use case logic.
func (m *MockResetPasswordUseCase) Execute(ctx context.Context, token, newPassword string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, token, newPassword)
	}
	return nil
}

// MockVerifyResetTokenUseCase is a mock implementation of VerifyResetTokenUseCase for testing.
type MockVerifyResetTokenUseCase struct {
	ExecuteFunc func(ctx context.Context, token string) error
}

// Execute executes the use case logic.
func (m *MockVerifyResetTokenUseCase) Execute(ctx context.Context, token string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, token)
	}
	return nil
}
