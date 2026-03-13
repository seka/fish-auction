package admin

import (
	"context"
	"errors"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type UpdatePasswordUseCase interface {
	Execute(ctx context.Context, id int, currentPassword, newPassword string) error
}

var _ UpdatePasswordUseCase = (*updatePasswordUseCase)(nil)

type updatePasswordUseCase struct {
	adminRepo repository.AdminRepository
}

// NewUpdatePasswordUseCase creates a new instance of UpdatePasswordUseCase
func NewUpdatePasswordUseCase(adminRepo repository.AdminRepository) *updatePasswordUseCase {
	return &updatePasswordUseCase{adminRepo: adminRepo}
}

func (uc *updatePasswordUseCase) Execute(ctx context.Context, id int, currentPassword, newPassword string) error {
	admin, err := uc.adminRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("invalid current password")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return uc.adminRepo.UpdatePassword(ctx, id, string(newHash))
}
