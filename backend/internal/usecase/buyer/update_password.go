package buyer

import (
	"context"
	"errors"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type UpdatePasswordUseCase interface {
	Execute(ctx context.Context, buyerID int, currentPassword, newPassword string) error
}

var _ UpdatePasswordUseCase = (*updatePasswordUseCase)(nil)

type updatePasswordUseCase struct {
	authRepo    repository.AuthenticationRepository
	sessionRepo repository.SessionRepository
}

// NewUpdatePasswordUseCase creates a new instance of UpdatePasswordUseCase
func NewUpdatePasswordUseCase(authRepo repository.AuthenticationRepository, sessionRepo repository.SessionRepository) *updatePasswordUseCase {
	return &updatePasswordUseCase{
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
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

	if err := uc.authRepo.UpdatePassword(ctx, buyerID, string(newHash)); err != nil {
		return err
	}

	// Invalidate all sessions after password change
	return uc.sessionRepo.DeleteAllByUserID(ctx, buyerID, model.SessionRoleBuyer)
}
