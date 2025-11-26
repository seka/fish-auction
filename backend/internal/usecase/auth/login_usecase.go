package auth

import (
	"context"
)

// LoginUseCase defines the interface for user authentication
type LoginUseCase interface {
	Execute(ctx context.Context, password string) (bool, error)
}

// loginUseCase handles user authentication
type loginUseCase struct{}

// NewLoginUseCase creates a new instance of LoginUseCase
func NewLoginUseCase() LoginUseCase {
	return &loginUseCase{}
}

// Execute authenticates a user with the provided password
func (uc *loginUseCase) Execute(ctx context.Context, password string) (bool, error) {
	if password == "admin-password" {
		return true, nil
	}
	return false, nil
}
