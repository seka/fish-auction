package buyer

import (
	"context"
	"errors"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type UpdatePasswordUseCase interface {
	Execute(ctx context.Context, buyerID int, currentPassword, newPassword string) error
}

type updatePasswordUseCase struct {
	authRepo repository.AuthenticationRepository
}

func NewUpdatePasswordUseCase(authRepo repository.AuthenticationRepository) UpdatePasswordUseCase {
	return &updatePasswordUseCase{authRepo: authRepo}
}

func (uc *updatePasswordUseCase) Execute(ctx context.Context, buyerID int, currentPassword, newPassword string) error {
	auth, err := uc.authRepo.FindByBuyerID(ctx, buyerID)
	if err != nil {
		return err
	}
	if auth == nil {
		return errors.New("buyer authentication not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("invalid current password")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return uc.authRepo.UpdatePassword(ctx, buyerID, string(newHash))
}
