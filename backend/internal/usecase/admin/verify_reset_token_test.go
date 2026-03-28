package admin_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
	"github.com/stretchr/testify/mock"
)

type mockPasswordResetRepositoryForVerify struct {
	mock.Mock
}

func (m *mockPasswordResetRepositoryForVerify) Create(ctx context.Context, userID int, role, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, role, tokenHash, expiresAt)
	return args.Error(0)
}

func (m *mockPasswordResetRepositoryForVerify) FindByTokenHash(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PasswordResetToken), args.Error(1)
}

func (m *mockPasswordResetRepositoryForVerify) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *mockPasswordResetRepositoryForVerify) DeleteAllByUserID(ctx context.Context, userID int, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func TestVerifyResetTokenUseCase_Execute(t *testing.T) {
	validToken := "valid-token"
	hash := sha256.Sum256([]byte(validToken))
	validTokenHash := hex.EncodeToString(hash[:])

	tests := []struct {
		name          string
		token         string
		mockTokenHash string
		mockRole      string
		mockExpiresAt time.Time
		wantErr       bool
	}{
		{
			name:          "Success",
			token:         validToken,
			mockTokenHash: validTokenHash,
			mockRole:      "admin",
			mockExpiresAt: time.Now().Add(1 * time.Hour),
			wantErr:       false,
		},
		{
			name:    "TokenNotFound",
			token:   "invalid-token",
			wantErr: true,
		},
		{
			name:          "WrongRole",
			token:         validToken,
			mockTokenHash: validTokenHash,
			mockRole:      "buyer",
			mockExpiresAt: time.Now().Add(1 * time.Hour),
			wantErr:       true,
		},
		{
			name:          "Expired",
			token:         validToken,
			mockTokenHash: validTokenHash,
			mockRole:      "admin",
			mockExpiresAt: time.Now().Add(-1 * time.Hour),
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockPasswordResetRepositoryForVerify{}

			hash := sha256.Sum256([]byte(tt.token))
			expectedHash := hex.EncodeToString(hash[:])

			if tt.mockTokenHash != "" {
				repo.On("FindByTokenHash", mock.Anything, expectedHash).Return(&model.PasswordResetToken{
					UserID:    1,
					Role:      tt.mockRole,
					ExpiresAt: tt.mockExpiresAt,
					TokenHash: expectedHash,
				}, nil)
				if tt.mockExpiresAt.Before(time.Now()) {
					repo.On("DeleteByTokenHash", mock.Anything, expectedHash).Return(nil)
				}
			} else {
				repo.On("FindByTokenHash", mock.Anything, expectedHash).Return(nil, nil)
			}

			uc := admin.NewVerifyResetTokenUseCase(repo)
			err := uc.Execute(context.Background(), tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}
