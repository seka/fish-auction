package buyer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepoForUpdate struct {
	auth      *model.Authentication
	findErr   error
	updateErr error
}

func (m *mockAuthRepoForUpdate) Login(_ context.Context, _, _ string) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockAuthRepoForUpdate) Create(_ context.Context, _ *model.Authentication) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepoForUpdate) FindByEmail(_ context.Context, _ string) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepoForUpdate) FindByBuyerID(_ context.Context, _ int) (*model.Authentication, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.auth, nil
}
func (m *mockAuthRepoForUpdate) UpdateLoginSuccess(_ context.Context, _ int, _ time.Time) error {
	return nil
}
func (m *mockAuthRepoForUpdate) IncrementFailedAttempts(_ context.Context, _ int) error {
	return nil
}
func (m *mockAuthRepoForUpdate) ResetFailedAttempts(_ context.Context, _ int) error {
	return nil
}
func (m *mockAuthRepoForUpdate) LockAccount(_ context.Context, _ int, _ time.Time) error {
	return nil
}
func (m *mockAuthRepoForUpdate) UpdatePassword(_ context.Context, _ int, _ string) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	return nil
}

type mockSessionRepo struct {
	repository.SessionRepository
}

func (m *mockSessionRepo) DeleteAllByUserID(_ context.Context, _ int, _ model.SessionRole) error {
	return nil
}

func TestUpdatePasswordUseCase_Execute(t *testing.T) {
	password := "current-password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	validAuth := &model.Authentication{BuyerID: 1, PasswordHash: string(hash)}

	tests := []struct {
		name        string
		buyerID     int
		currentPass string
		newPass     string
		mockAuth    *model.Authentication
		findErr     error
		updateErr   error
		wantErr     bool
	}{
		{
			name:        "Success",
			buyerID:     1,
			currentPass: "current-password",
			newPass:     "new-password",
			mockAuth:    validAuth,
		},
		{
			name:        "IncorrectCurrentPassword",
			buyerID:     1,
			currentPass: "wrong-password",
			newPass:     "new-password",
			mockAuth:    validAuth,
			wantErr:     true,
		},
		{
			name:        "NotFound",
			buyerID:     99,
			currentPass: "current-password",
			mockAuth:    nil,
			wantErr:     true,
		},
		{
			name:        "FindRepoError",
			buyerID:     1,
			currentPass: "current-password",
			mockAuth:    validAuth,
			findErr:     errors.New("find error"),
			wantErr:     true,
		},
		{
			name:        "UpdateRepoError",
			buyerID:     1,
			currentPass: "current-password",
			newPass:     "new-password",
			mockAuth:    validAuth,
			updateErr:   errors.New("update error"),
			wantErr:     true,
		},
		{
			name:        "PasswordTooLong",
			buyerID:     1,
			currentPass: "current-password",
			newPass:     "this_password_is_definitely_way_too_long_to_be_hashed_by_bcrypt_because_it_exceeds_seventy_two_bytes_limit",
			mockAuth:    validAuth,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAuthRepoForUpdate{
				auth:      tt.mockAuth,
				findErr:   tt.findErr,
				updateErr: tt.updateErr,
			}
			uc := buyer.NewUpdatePasswordUseCase(repo, &mockSessionRepo{})

			err := uc.Execute(context.Background(), tt.buyerID, tt.currentPass, tt.newPass)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
