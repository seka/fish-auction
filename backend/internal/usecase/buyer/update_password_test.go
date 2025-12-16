package buyer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepoForUpdate struct {
	auth *model.Authentication
	err  error
}

func (m *mockAuthRepoForUpdate) Login(ctx context.Context, email, password string) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockAuthRepoForUpdate) Create(ctx context.Context, auth *model.Authentication) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepoForUpdate) FindByEmail(ctx context.Context, email string) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepoForUpdate) FindByBuyerID(ctx context.Context, buyerID int) (*model.Authentication, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.auth, nil
}
func (m *mockAuthRepoForUpdate) UpdateLoginSuccess(ctx context.Context, id int, loginAt time.Time) error {
	return nil
}
func (m *mockAuthRepoForUpdate) IncrementFailedAttempts(ctx context.Context, id int) error {
	return nil
}
func (m *mockAuthRepoForUpdate) ResetFailedAttempts(ctx context.Context, id int) error {
	return nil
}
func (m *mockAuthRepoForUpdate) LockAccount(ctx context.Context, id int, until time.Time) error {
	return nil
}
func (m *mockAuthRepoForUpdate) UpdatePassword(ctx context.Context, buyerID int, hashedPassword string) error {
	if m.err != nil {
		return m.err
	}
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
		repoErr     error
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
			name:        "RepoError",
			buyerID:     1,
			currentPass: "current-password",
			mockAuth:    validAuth,
			repoErr:     errors.New("db error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAuthRepoForUpdate{auth: tt.mockAuth, err: tt.repoErr}
			uc := buyer.NewUpdatePasswordUseCase(repo)

			err := uc.Execute(context.Background(), tt.buyerID, tt.currentPass, tt.newPass)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
