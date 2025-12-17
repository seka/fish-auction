package admin_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
	"github.com/stretchr/testify/mock"
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

type mockAdminPasswordResetRepositoryForReset struct {
	mock.Mock
}

func (m *mockAdminPasswordResetRepositoryForReset) Create(ctx context.Context, userID int, role string, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, role, tokenHash, expiresAt)
	return args.Error(0)
}

func (m *mockAdminPasswordResetRepositoryForReset) FindByTokenHash(ctx context.Context, tokenHash string) (int, string, time.Time, error) {
	args := m.Called(ctx, tokenHash)
	return args.Int(0), args.String(1), args.Get(2).(time.Time), args.Error(3)
}

func (m *mockAdminPasswordResetRepositoryForReset) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *mockAdminPasswordResetRepositoryForReset) DeleteAllByUserID(ctx context.Context, userID int, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func TestResetPasswordUseCase_Execute(t *testing.T) {
	validAdminID := 1
	validExpires := time.Now().Add(1 * time.Hour)
	validToken := "valid-token"
	hash := sha256.Sum256([]byte(validToken))
	validTokenHash := hex.EncodeToString(hash[:])

	tests := []struct {
		name          string
		token         string
		newPassword   string
		mockTokenHash string
		mockAdminID   int
		mockExpiresAt time.Time
		mockFindErr   error
		mockUpdateErr error
		wantErr       bool
	}{
		{
			name:          "Success",
			token:         validToken,
			newPassword:   "newPass123",
			mockTokenHash: validTokenHash,
			mockAdminID:   validAdminID,
			mockExpiresAt: validExpires,
			wantErr:       false,
		},
		{
			name:        "TokenNotFound",
			token:       "invalid",
			newPassword: "newPass123",
			wantErr:     true,
		},
		{
			name:        "TokenFindError",
			token:       validToken,
			newPassword: "newPass123",
			mockFindErr: errors.New("db error"),
			wantErr:     true,
		},
		{
			name:          "TokenExpired",
			token:         validToken,
			newPassword:   "newPass123",
			mockTokenHash: validTokenHash,
			mockAdminID:   validAdminID,
			mockExpiresAt: time.Now().Add(-1 * time.Hour),
			wantErr:       true,
		},
		{
			name:          "UpdateFailed",
			token:         validToken,
			newPassword:   "newPass123",
			mockTokenHash: validTokenHash,
			mockAdminID:   validAdminID,
			mockExpiresAt: validExpires,
			mockUpdateErr: errors.New("update failed"),
			wantErr:       true,
		},
		{
			name:          "PasswordTooLong",
			token:         validToken,
			newPassword:   "this_password_is_definitely_way_too_long_to_be_hashed_by_bcrypt_because_it_exceeds_seventy_two_bytes_limit",
			mockTokenHash: validTokenHash,
			mockAdminID:   validAdminID,
			mockExpiresAt: validExpires,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwdResetRepo := &mockAdminPasswordResetRepositoryForReset{}

			// Calculate expected hash for the input token
			hash := sha256.Sum256([]byte(tt.token))
			expectedHash := hex.EncodeToString(hash[:])

			if tt.mockFindErr != nil {
				pwdResetRepo.On("FindByTokenHash", mock.Anything, expectedHash).Return(0, "", time.Time{}, tt.mockFindErr)
			} else if tt.mockTokenHash == "" {
				// Token not found scenario
				pwdResetRepo.On("FindByTokenHash", mock.Anything, expectedHash).Return(0, "", time.Time{}, nil)
			} else {
				// Found
				// Ensure that the mockTokenHash in test case matches expectedHash?
				// Actually for the "Success" case, tt.mockTokenHash is validTokenHash, and tt.token is validToken. So they match.
				// For "TokenExpired", same.
				pwdResetRepo.On("FindByTokenHash", mock.Anything, expectedHash).Return(tt.mockAdminID, "admin", tt.mockExpiresAt, nil)

				if tt.mockAdminID != 0 {
					if tt.mockExpiresAt.After(time.Now()) {
						// Valid
						pwdResetRepo.On("DeleteAllByUserID", mock.Anything, tt.mockAdminID, "admin").Return(nil)
					} else {
						// Expired
						pwdResetRepo.On("DeleteByTokenHash", mock.Anything, expectedHash).Return(nil)
					}
				}
			}

			adminRepo := &mockAdminRepositoryForReset{err: tt.mockUpdateErr}

			uc := admin.NewResetPasswordUseCase(pwdResetRepo, adminRepo)
			err := uc.Execute(context.Background(), tt.token, tt.newPassword)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}
