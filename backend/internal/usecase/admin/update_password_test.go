package admin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
	"golang.org/x/crypto/bcrypt"
)

// Mock for UpdatePassword
type mockAdminRepositoryForUpdate struct {
	admin     *entity.Admin
	findErr   error
	updateErr error
}

func (m *mockAdminRepositoryForUpdate) FindByID(ctx context.Context, id int) (*entity.Admin, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	if m.admin != nil && m.admin.ID == id {
		return m.admin, nil
	}
	return nil, nil
}
func (m *mockAdminRepositoryForUpdate) FindOneByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	return nil, nil
}
func (m *mockAdminRepositoryForUpdate) Create(ctx context.Context, admin *entity.Admin) error {
	return nil
}
func (m *mockAdminRepositoryForUpdate) Count(ctx context.Context) (int, error) {
	return 0, nil
}
func (m *mockAdminRepositoryForUpdate) UpdatePassword(ctx context.Context, id int, hash string) error {
	return m.updateErr
}

func TestUpdatePasswordUseCase_Execute(t *testing.T) {
	// Prepare a hashed password
	password := "correctPassword"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	adminUser := &entity.Admin{ID: 1, PasswordHash: string(hashed)}

	tests := []struct {
		name            string
		id              int
		currentPassword string
		newPassword     string
		mockAdmin       *entity.Admin
		findErr         error
		updateErr       error
		wantErr         bool
	}{
		{
			name:            "Success",
			id:              1,
			currentPassword: password,
			newPassword:     "newPass123",
			mockAdmin:       adminUser,
			wantErr:         false,
		},
		{
			name:            "AdminNotFound",
			id:              99,
			currentPassword: password,
			newPassword:     "newPass123",
			mockAdmin:       adminUser, // Returns nil for ID 99
			wantErr:         true,
		},
		{
			name:            "IncorrectCurrentPassword",
			id:              1,
			currentPassword: "wrongPassword",
			newPassword:     "newPass123",
			mockAdmin:       adminUser,
			wantErr:         true,
		},
		{
			name:            "RepoFindError",
			id:              1,
			currentPassword: password,
			newPassword:     "newPass123",
			findErr:         errors.New("db error"),
			wantErr:         true,
		},
		{
			name:            "RepoUpdateError",
			id:              1,
			currentPassword: password,
			newPassword:     "newPass123",
			mockAdmin:       adminUser,
			updateErr:       errors.New("update failed"),
			wantErr:         true,
		},
		{
			name:            "NewPasswordTooLong",
			id:              1,
			currentPassword: password,
			newPassword:     "this_password_is_definitely_way_too_long_to_be_hashed_by_bcrypt_because_it_exceeds_seventy_two_bytes_limit",
			mockAdmin:       adminUser,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAdminRepositoryForUpdate{
				admin:     tt.mockAdmin,
				findErr:   tt.findErr,
				updateErr: tt.updateErr,
			}
			uc := admin.NewUpdatePasswordUseCase(repo)
			err := uc.Execute(context.Background(), tt.id, tt.currentPassword, tt.newPassword)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
