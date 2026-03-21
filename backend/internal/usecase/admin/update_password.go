package admin

import (
	"context"
	"errors"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// UpdatePasswordUseCase defines the interface for updating an admin password.
type UpdatePasswordUseCase interface {
	// Execute updates the admin password for the given ID.
	Execute(ctx context.Context, id int, currentPassword, newPassword string) error
}

var _ UpdatePasswordUseCase = (*updatePasswordUseCase)(nil)

type updatePasswordUseCase struct {
	adminRepo   repository.AdminRepository
	sessionRepo repository.SessionRepository
}

// NewUpdatePasswordUseCase creates a new UpdatePasswordUseCase instance.
func NewUpdatePasswordUseCase(adminRepo repository.AdminRepository, sessionRepo repository.SessionRepository) *updatePasswordUseCase {
	return &updatePasswordUseCase{
		adminRepo:   adminRepo,
		sessionRepo: sessionRepo,
	}
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

	if err := uc.adminRepo.UpdatePassword(ctx, id, string(newHash)); err != nil {
		return err
	}

	// Invalidate all sessions after password change
	return uc.sessionRepo.DeleteAllByUserID(ctx, id, model.SessionRoleAdmin)
}
