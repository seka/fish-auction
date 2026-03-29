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
	usetesting "github.com/seka/fish-auction/backend/internal/usecase/testing"
	"github.com/stretchr/testify/mock"
)

type mockAuthRepositoryForReset struct {
	err error
}

func (m *mockAuthRepositoryForReset) Login(_ context.Context, _, _ string) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockAuthRepositoryForReset) Create(_ context.Context, _ *model.Authentication) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepositoryForReset) FindByEmail(_ context.Context, _ string) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepositoryForReset) FindByBuyerID(_ context.Context, _ int) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepositoryForReset) UpdateLoginSuccess(_ context.Context, _ int, _ time.Time) error {
	return nil
}
func (m *mockAuthRepositoryForReset) IncrementFailedAttempts(_ context.Context, _ int) error {
	return nil
}
func (m *mockAuthRepositoryForReset) ResetFailedAttempts(_ context.Context, _ int) error {
	return nil
}
func (m *mockAuthRepositoryForReset) LockAccount(_ context.Context, _ int, _ time.Time) error {
	return nil
}
func (m *mockAuthRepositoryForReset) UpdatePassword(_ context.Context, _ int, _ string) error {
	return m.err
}

type mockBuyerPasswordResetRepositoryForReset struct {
	mock.Mock
}

func (m *mockBuyerPasswordResetRepositoryForReset) Create(ctx context.Context, userID int, role, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, role, tokenHash, expiresAt)
	return args.Error(0)
}
func (m *mockBuyerPasswordResetRepositoryForReset) FindByTokenHash(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PasswordResetToken), args.Error(1)
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
			newPassword:   "NewPassword123!",
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
			newPassword:   "NewPassword123!",
			mockTokenHash: validTokenHash,
			mockBuyerID:   validBuyerID,
			mockExpiresAt: validExpires,
			mockUpdateErr: errors.New("update failed"),
			wantErr:       true,
		},
		{
			name:          "RoleMismatch",
			token:         validToken,
			newPassword:   "NewPassword123!",
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

			switch {
			case tt.mockFindErr != nil:
				resetRepo.On("FindByTokenHash", mock.Anything, expectedHash).Return(nil, tt.mockFindErr)
			case tt.mockTokenHash == "": // Token not found scenario
				resetRepo.On("FindByTokenHash", mock.Anything, expectedHash).Return(nil, nil)
			default: // Token found scenario
				role := "buyer"
				if tt.name == "RoleMismatch" {
					role = "admin" // Wrong role for buyer usecase
				}
				resetModel := &model.PasswordResetToken{
					UserID:    tt.mockBuyerID,
					Role:      role,
					ExpiresAt: tt.mockExpiresAt,
					TokenHash: expectedHash,
				}
				resetRepo.On("FindByTokenHash", mock.Anything, expectedHash).Return(resetModel, nil)
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
			txMgr := &usetesting.MockTransactionManager{}

			uc := auth.NewResetPasswordUseCase(resetRepo, authRepo, txMgr)
			err := uc.Execute(context.Background(), tt.token, tt.newPassword)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}
