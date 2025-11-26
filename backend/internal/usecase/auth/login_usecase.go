package auth

import (
	"context"
)

// LoginUseCase handles user authentication
type LoginUseCase struct{}

// NewLoginUseCase creates a new instance of LoginUseCase
func NewLoginUseCase() *LoginUseCase {
	return &LoginUseCase{}
}

// Execute authenticates a user with the provided password
func (uc *LoginUseCase) Execute(ctx context.Context, password string) (bool, error) {
	if password == "admin-password" {
		return true, nil
	}
	return false, nil
}
