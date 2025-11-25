package usecase

import (
	"context"
)

// AuthUseCase defines the interface for authentication-related business logic
type AuthUseCase interface {
	Login(ctx context.Context, password string) (bool, error)
}

type authInteractor struct{}

func NewAuthInteractor() AuthUseCase {
	return &authInteractor{}
}

func (i *authInteractor) Login(ctx context.Context, password string) (bool, error) {
	if password == "admin-password" {
		return true, nil
	}
	return false, nil
}
