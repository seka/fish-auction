package admin_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
)

type mockAdminRepositoryForReset struct {
	err error
}

func (m *mockAdminRepositoryForReset) FindByID(ctx context.Context, id int) (*entity.Admin, error) {
	return nil, nil
}
func (m *mockAdminRepositoryForReset) FindOneByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	return nil, nil
}
func (m *mockAdminRepositoryForReset) Create(ctx context.Context, admin *entity.Admin) error {
	return nil
}
func (m *mockAdminRepositoryForReset) Count(ctx context.Context) (int, error) { return 0, nil }
func (m *mockAdminRepositoryForReset) UpdatePassword(ctx context.Context, id int, hash string) error {
	return m.err
}

func TestResetPasswordUseCase_Execute(t *testing.T) {
	validAdminID := 1
	validExpires := time.Now().Add(1 * time.Hour)
	expiredExpires := time.Now().Add(-1 * time.Hour)

	tests := []struct {
		name        string
		token       string
		newPassword string
		mockAdminID int
		mockExpires time.Time
		repoFound   bool
		updateErr   error
		wantError   bool
	}{
		{
			name:        "Success",
			token:       "valid-token",
			newPassword: "new-password",
			mockAdminID: validAdminID,
			mockExpires: validExpires,
			repoFound:   true,
		},
		{
			name:      "TokenNotFound",
			token:     "invalid-token",
			repoFound: false,
			wantError: true,
		},
		{
			name:        "TokenExpired",
			token:       "expired-token",
			mockAdminID: validAdminID,
			mockExpires: expiredExpires,
			repoFound:   true,
			wantError:   true,
		},
		{
			name:        "UpdateFailed",
			token:       "valid-token",
			newPassword: "new-password",
			mockAdminID: validAdminID,
			mockExpires: validExpires,
			repoFound:   true,
			updateErr:   errors.New("update failed"),
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRepo := &mockAdminPasswordResetRepositoryWithReturn{
				adminID:   tt.mockAdminID,
				expiresAt: tt.mockExpires,
				found:     tt.repoFound,
			}
			adminRepo := &mockAdminRepositoryForReset{err: tt.updateErr}

			uc := admin.NewResetPasswordUseCase(resetRepo, adminRepo)
			err := uc.Execute(context.Background(), tt.token, tt.newPassword)

			if (err != nil) != tt.wantError {
				t.Errorf("expected error=%v, got %v", tt.wantError, err)
			}
		})
	}
}

type mockAdminPasswordResetRepositoryWithReturn struct {
	adminID   int
	expiresAt time.Time
	found     bool
}

func (m *mockAdminPasswordResetRepositoryWithReturn) Create(ctx context.Context, adminID int, tokenHash string, expiresAt time.Time) error {
	return nil
}
func (m *mockAdminPasswordResetRepositoryWithReturn) FindByTokenHash(ctx context.Context, tokenHash string) (int, time.Time, error) {
	if !m.found {
		return 0, time.Time{}, errors.New("not found")
	}
	return m.adminID, m.expiresAt, nil
}
func (m *mockAdminPasswordResetRepositoryWithReturn) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	return nil
}
func (m *mockAdminPasswordResetRepositoryWithReturn) DeleteAllByAdminID(ctx context.Context, adminID int) error {
	return nil
}
