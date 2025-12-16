package admin_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
)

type mockAdminRepoForReqPwd struct {
	admin *entity.Admin
	err   error
}

func (m *mockAdminRepoForReqPwd) Create(ctx context.Context, admin *entity.Admin) error { return nil }
func (m *mockAdminRepoForReqPwd) FindOneByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.admin != nil && m.admin.Email == email {
		return m.admin, nil
	}
	return nil, nil
}
func (m *mockAdminRepoForReqPwd) FindByID(ctx context.Context, id int) (*entity.Admin, error) {
	return nil, nil
}
func (m *mockAdminRepoForReqPwd) UpdatePassword(ctx context.Context, id int, password string) error {
	return nil
}
func (m *mockAdminRepoForReqPwd) Count(ctx context.Context) (int, error) { return 0, nil }

type mockPwdResetRepoForReqPwd struct {
	delErr error
	creErr error
}

func (m *mockPwdResetRepoForReqPwd) Create(ctx context.Context, adminID int, tokenHash string, expiresAt time.Time) error {
	return m.creErr
}
func (m *mockPwdResetRepoForReqPwd) FindByTokenHash(ctx context.Context, tokenHash string) (int, time.Time, error) {
	return 0, time.Time{}, nil
}
func (m *mockPwdResetRepoForReqPwd) DeleteAllByAdminID(ctx context.Context, adminID int) error {
	return m.delErr
}
func (m *mockPwdResetRepoForReqPwd) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	return nil
}

type mockEmailServiceForReqPwd struct {
	sndErr error
}

func (m *mockEmailServiceForReqPwd) SendAdminPasswordReset(ctx context.Context, to, url string) error {
	return m.sndErr
}
func (m *mockEmailServiceForReqPwd) SendBuyerPasswordReset(ctx context.Context, to, url string) error {
	return nil
}

func TestRequestPasswordResetUseCase_Execute(t *testing.T) {
	validAdmin := &entity.Admin{ID: 1, Email: "test@example.com"}

	tests := []struct {
		name        string
		email       string
		mockAdmin   *entity.Admin
		mockFindErr error
		mockRandErr error
		mockDelErr  error
		mockCreErr  error
		mockSndErr  error
		wantErr     bool
	}{
		{
			name:      "Success",
			email:     "test@example.com",
			mockAdmin: validAdmin,
		},
		{
			name:      "AdminNotFound",
			email:     "unknown@example.com",
			mockAdmin: nil,   // Repo returns nil, nil
			wantErr:   false, // Should return nil (masking)
		},
		{
			name:        "RepoFindError",
			email:       "test@example.com",
			mockFindErr: errors.New("db error"),
			wantErr:     false, // Should return nil (masking)
		},
		{
			name:        "RandReadError",
			email:       "test@example.com",
			mockAdmin:   validAdmin,
			mockRandErr: errors.New("rand error"),
			wantErr:     true,
		},
		{
			name:       "DeleteTokenError",
			email:      "test@example.com",
			mockAdmin:  validAdmin,
			mockDelErr: errors.New("del error"),
			wantErr:    true,
		},
		{
			name:       "CreateTokenError",
			email:      "test@example.com",
			mockAdmin:  validAdmin,
			mockCreErr: errors.New("cre error"),
			wantErr:    true,
		},
		{
			name:       "SendEmailError",
			email:      "test@example.com",
			mockAdmin:  validAdmin,
			mockSndErr: errors.New("send error"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock rand
			restore := admin.SetRandRead(admin.GetRandReadFunc(tt.mockRandErr))
			defer restore()

			adminRepo := &mockAdminRepoForReqPwd{admin: tt.mockAdmin, err: tt.mockFindErr}
			pwdResetRepo := &mockPwdResetRepoForReqPwd{delErr: tt.mockDelErr, creErr: tt.mockCreErr}
			emailService := &mockEmailServiceForReqPwd{sndErr: tt.mockSndErr}

			uc := admin.NewRequestPasswordResetUseCase(adminRepo, pwdResetRepo, emailService)
			err := uc.Execute(context.Background(), tt.email)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
