package admin_test

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
	usetesting "github.com/seka/fish-auction/backend/internal/usecase/testing"
	"github.com/stretchr/testify/mock"
)

type mockAdminRepoForReqPwd struct {
	admin *model.Admin
	err   error
}

func (m *mockAdminRepoForReqPwd) Create(_ context.Context, _ *model.Admin) error { return nil }
func (m *mockAdminRepoForReqPwd) FindOneByEmail(_ context.Context, email string) (*model.Admin, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.admin != nil && m.admin.Email == email {
		return m.admin, nil
	}
	return nil, nil
}
func (m *mockAdminRepoForReqPwd) FindByID(_ context.Context, _ int) (*model.Admin, error) {
	return nil, nil
}
func (m *mockAdminRepoForReqPwd) UpdatePassword(_ context.Context, _ int, _ string) error {
	return nil
}
func (m *mockAdminRepoForReqPwd) Count(_ context.Context) (int, error) { return 0, nil }

type mockPwdResetRepoForReqPwd struct {
	mock.Mock
}

func (m *mockPwdResetRepoForReqPwd) Create(ctx context.Context, userID int, role, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, role, tokenHash, expiresAt)
	return args.Error(0)
}
func (m *mockPwdResetRepoForReqPwd) FindByTokenHash(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PasswordResetToken), args.Error(1)
}
func (m *mockPwdResetRepoForReqPwd) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}
func (m *mockPwdResetRepoForReqPwd) DeleteAllByUserID(ctx context.Context, userID int, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

type mockEmailServiceForReqPwd struct {
	sndErr error
}

func (m *mockEmailServiceForReqPwd) SendAdminPasswordReset(_ context.Context, _, _ string) error {
	return m.sndErr
}
func (m *mockEmailServiceForReqPwd) SendBuyerPasswordReset(_ context.Context, _, _ string) error {
	return nil
}

func TestRequestPasswordResetUseCase_Execute(t *testing.T) {
	validAdmin := &model.Admin{ID: 1, Email: "test@example.com"}

	tests := []struct {
		name        string
		email       string
		mockAdmin   *model.Admin
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
			wantErr:   true,
		},
		{
			name:        "RepoFindError",
			email:       "test@example.com",
			mockFindErr: errors.New("db error"),
			wantErr:     true,
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
			pwdResetRepo := &mockPwdResetRepoForReqPwd{}

			if tt.mockFindErr == nil && tt.mockAdmin != nil {
				// Happy path for finding admin.
				pwdResetRepo.On("DeleteAllByUserID", mock.Anything, tt.mockAdmin.ID, "admin").Return(tt.mockDelErr)
				if tt.mockDelErr == nil {
					pwdResetRepo.On("Create", mock.Anything, tt.mockAdmin.ID, "admin", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(tt.mockCreErr)
				}
			}

			emailService := &mockEmailServiceForReqPwd{sndErr: tt.mockSndErr}
			txMgr := &usetesting.MockTransactionManager{}

			frontendURL, _ := url.Parse("https://localhost")
			uc := admin.NewRequestPasswordResetUseCase(adminRepo, pwdResetRepo, emailService, frontendURL, txMgr)
			err := uc.Execute(context.Background(), tt.email)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
