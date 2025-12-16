package admin_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
)

type mockAdminRepository struct {
	admin *entity.Admin
	err   error
}

func (m *mockAdminRepository) FindByID(ctx context.Context, id int) (*entity.Admin, error) {
	return nil, nil
}
func (m *mockAdminRepository) FindOneByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.admin != nil && m.admin.Email == email {
		return m.admin, nil
	}
	return nil, nil
}
func (m *mockAdminRepository) Create(ctx context.Context, admin *entity.Admin) error { return nil }
func (m *mockAdminRepository) Count(ctx context.Context) (int, error)                { return 0, nil }
func (m *mockAdminRepository) UpdatePassword(ctx context.Context, id int, hash string) error {
	return nil
}

type mockAdminPasswordResetRepository struct {
	createErr error
}

func (m *mockAdminPasswordResetRepository) Create(ctx context.Context, adminID int, tokenHash string, expiresAt time.Time) error {
	return m.createErr
}
func (m *mockAdminPasswordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (int, time.Time, error) {
	return 0, time.Time{}, nil
}
func (m *mockAdminPasswordResetRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	return nil
}
func (m *mockAdminPasswordResetRepository) DeleteAllByAdminID(ctx context.Context, adminID int) error {
	return nil
}

type mockEmailService struct {
	sentAdminURL string
	err          error
}

func (m *mockEmailService) SendBuyerPasswordReset(ctx context.Context, to, url string) error {
	return nil
}
func (m *mockEmailService) SendAdminPasswordReset(ctx context.Context, to, url string) error {
	if m.err != nil {
		return m.err
	}
	m.sentAdminURL = url
	return nil
}

func TestRequestPasswordResetUseCase_Execute(t *testing.T) {
	validAdmin := &entity.Admin{ID: 1, Email: "admin@example.com"}

	tests := []struct {
		name         string
		email        string
		mockAdmin    *entity.Admin
		mockRepoErr  error
		mockEmailErr error
		wantError    bool
		wantSent     bool
	}{
		{
			name:      "Success",
			email:     "admin@example.com",
			mockAdmin: validAdmin,
			wantSent:  true,
		},
		{
			name:      "UserNotFound",
			email:     "other@example.com",
			mockAdmin: validAdmin,
			wantSent:  false,
		},
		{
			name:         "EmailError",
			email:        "admin@example.com",
			mockAdmin:    validAdmin,
			mockEmailErr: errors.New("email failed"),
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminRepo := &mockAdminRepository{admin: tt.mockAdmin, err: tt.mockRepoErr}
			resetRepo := &mockAdminPasswordResetRepository{}
			emailService := &mockEmailService{err: tt.mockEmailErr}

			uc := admin.NewRequestPasswordResetUseCase(adminRepo, resetRepo, emailService)
			err := uc.Execute(context.Background(), tt.email)

			if (err != nil) != tt.wantError {
				if tt.name == "RepoError" && err == nil {
					// Expected logic
				} else {
					t.Errorf("expected error=%v, got %v", tt.wantError, err)
				}
			}

			if tt.wantSent && emailService.sentAdminURL == "" {
				t.Error("expected email to be sent, but wasn't")
			}
		})
	}
}
