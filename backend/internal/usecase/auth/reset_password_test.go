package auth_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
	"github.com/stretchr/testify/mock"
)

type mockAuthRepositoryForReset struct {
	err error
}

func (m *mockAuthRepositoryForReset) Login(ctx context.Context, email, password string) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockAuthRepositoryForReset) Create(ctx context.Context, auth *model.Authentication) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepositoryForReset) FindByEmail(ctx context.Context, email string) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepositoryForReset) FindByBuyerID(ctx context.Context, buyerID int) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepositoryForReset) UpdateLoginSuccess(ctx context.Context, id int, loginAt time.Time) error {
	return nil
}
func (m *mockAuthRepositoryForReset) IncrementFailedAttempts(ctx context.Context, id int) error {
	return nil
}
func (m *mockAuthRepositoryForReset) ResetFailedAttempts(ctx context.Context, id int) error {
	return nil
}
func (m *mockAuthRepositoryForReset) LockAccount(ctx context.Context, id int, until time.Time) error {
	return nil
}
func (m *mockAuthRepositoryForReset) UpdatePassword(ctx context.Context, buyerID int, hashedPassword string) error {
	return m.err
}

type mockBuyerPasswordResetRepositoryForReset struct {
	mock.Mock
}

func (m *mockBuyerPasswordResetRepositoryForReset) Create(ctx context.Context, userID int, role string, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, role, tokenHash, expiresAt)
	return args.Error(0)
}
func (m *mockBuyerPasswordResetRepositoryForReset) FindByTokenHash(ctx context.Context, tokenHash string) (int, string, time.Time, error) {
	args := m.Called(ctx, tokenHash)
	return args.Int(0), args.String(1), args.Get(2).(time.Time), args.Error(3)
}
func (m *mockBuyerPasswordResetRepositoryForReset) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}
func (m *mockBuyerPasswordResetRepositoryForReset) DeleteAllByUserID(ctx context.Context, userID int, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func TestResetPasswordUseCase_Execute(t *testing.T) {
	validBuyerID := 1
	validExpires := time.Now().Add(1 * time.Hour)
	validToken := "valid-token"

	hash := sha256.Sum256([]byte(validToken))
	validTokenHash := hex.EncodeToString(hash[:])

	tests := []struct {
		name          string
		token         string
		newPassword   string
		mockTokenHash string
		mockBuyerID   int
		mockExpiresAt time.Time
		mockFindErr   error
		mockUpdateErr error
		wantErr       bool
	}{
		{
			name:          "Success",
			token:         validToken,
			newPassword:   "new-password",
			mockTokenHash: validTokenHash,
			mockBuyerID:   validBuyerID,
			mockExpiresAt: validExpires,
			wantErr:       false,
		},
		{
			name:        "TokenNotFound",
			token:       "invalid-token",
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
			mockBuyerID:   validBuyerID,
			mockExpiresAt: time.Now().Add(-1 * time.Hour),
			wantErr:       true,
		},
		{
			name:          "UpdateFailed",
			token:         validToken,
			newPassword:   "new-password",
			mockTokenHash: validTokenHash,
			mockBuyerID:   validBuyerID,
			mockExpiresAt: validExpires,
			mockUpdateErr: errors.New("update failed"),
			wantErr:       true,
		},
		{
			name:          "PasswordTooLong",
			token:         validToken,
			newPassword:   "this_password_is_definitely_way_too_long_to_be_hashed_by_bcrypt_because_it_exceeds_seventy_two_bytes_limit",
			mockTokenHash: validTokenHash,
			mockBuyerID:   validBuyerID,
			mockExpiresAt: validExpires,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRepo := &mockBuyerPasswordResetRepositoryForReset{}

			hash := sha256.Sum256([]byte(tt.token))
			expectedHash := hex.EncodeToString(hash[:])

			if tt.mockFindErr != nil {
				resetRepo.On("FindByTokenHash", mock.Anything, expectedHash).Return(0, "", time.Time{}, tt.mockFindErr)
			} else if tt.mockTokenHash == "" {
				resetRepo.On("FindByTokenHash", mock.Anything, expectedHash).Return(0, "", time.Time{}, nil)
			} else {
				resetRepo.On("FindByTokenHash", mock.Anything, expectedHash).Return(tt.mockBuyerID, "buyer", tt.mockExpiresAt, nil)
				if tt.mockBuyerID != 0 {
					if tt.mockExpiresAt.After(time.Now()) {
						// Valid
						resetRepo.On("DeleteAllByUserID", mock.Anything, tt.mockBuyerID, "buyer").Return(nil)
					} else {
						// Expired
						resetRepo.On("DeleteByTokenHash", mock.Anything, expectedHash).Return(nil)
					}
				}
			}
			authRepo := &mockAuthRepositoryForReset{err: tt.mockUpdateErr}

			uc := auth.NewResetPasswordUseCase(resetRepo, authRepo)
			err := uc.Execute(context.Background(), tt.token, tt.newPassword)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}
