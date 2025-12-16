package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
)

type mockAuthRepository struct {
	err error
}

func (m *mockAuthRepository) Login(ctx context.Context, email, password string) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockAuthRepository) Create(ctx context.Context, auth *model.Authentication) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepository) FindByEmail(ctx context.Context, email string) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepository) FindByBuyerID(ctx context.Context, buyerID int) (*model.Authentication, error) {
	return nil, nil
}
func (m *mockAuthRepository) UpdateLoginSuccess(ctx context.Context, id int, loginAt time.Time) error {
	return nil
}
func (m *mockAuthRepository) IncrementFailedAttempts(ctx context.Context, id int) error {
	return nil
}
func (m *mockAuthRepository) ResetFailedAttempts(ctx context.Context, id int) error {
	return nil
}
func (m *mockAuthRepository) LockAccount(ctx context.Context, id int, until time.Time) error {
	return nil
}
func (m *mockAuthRepository) UpdatePassword(ctx context.Context, buyerID int, hashedPassword string) error {
	return m.err
}

func TestResetPasswordUseCase_Execute(t *testing.T) {
	// For Reset logic, repo returns (buyerID, expiresAt, error)
	validBuyerID := 1
	validExpires := time.Now().Add(1 * time.Hour)
	expiredExpires := time.Now().Add(-1 * time.Hour)

	tests := []struct {
		name        string
		token       string
		newPassword string
		mockBuyerID int
		mockExpires time.Time
		repoFound   bool
		repoErr     error
		updateErr   error
		wantError   bool
	}{
		{
			name:        "Success",
			token:       "valid-token",
			newPassword: "new-password",
			mockBuyerID: validBuyerID,
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
			mockBuyerID: validBuyerID,
			mockExpires: expiredExpires,
			repoFound:   true,
			wantError:   true,
		},
		{
			name:        "UpdateFailed",
			token:       "valid-token",
			newPassword: "new-password",
			mockBuyerID: validBuyerID,
			mockExpires: validExpires,
			repoFound:   true,
			updateErr:   errors.New("update failed"),
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRepo := &mockBuyerPasswordResetRepositoryWithReturn{
				buyerID:   tt.mockBuyerID,
				expiresAt: tt.mockExpires,
				found:     tt.repoFound,
				err:       tt.repoErr,
			}
			authRepo := &mockAuthRepository{err: tt.updateErr}

			uc := auth.NewResetPasswordUseCase(resetRepo, authRepo)
			err := uc.Execute(context.Background(), tt.token, tt.newPassword)

			if (err != nil) != tt.wantError {
				t.Errorf("expected error=%v, got %v", tt.wantError, err)
			}
		})
	}
}

// Special mock for this test to return specific token
type mockBuyerPasswordResetRepositoryWithReturn struct {
	buyerID   int
	expiresAt time.Time
	found     bool
	err       error
}

func (m *mockBuyerPasswordResetRepositoryWithReturn) Create(ctx context.Context, buyerID int, tokenHash string, expiresAt time.Time) error {
	return nil
}
func (m *mockBuyerPasswordResetRepositoryWithReturn) FindByTokenHash(ctx context.Context, tokenHash string) (int, time.Time, error) {
	if m.err != nil {
		return 0, time.Time{}, m.err
	}
	if !m.found {
		return 0, time.Time{}, errors.New("not found") // UseCase calls repo.FindByTokenHash, if err != nil returns err.
		// Actually typical repo returns specific NotFound error or just error.
		// If usecase relies on error for not found, then this is fine.
		// If usecase expects (0, zeroTime, nil), logic is different.
		// Logic: buyerID, expiresAt, err := u.pwdResetRepo.FindByTokenHash... if err != nil return err.
	}
	return m.buyerID, m.expiresAt, nil
}
func (m *mockBuyerPasswordResetRepositoryWithReturn) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	return nil
}
func (m *mockBuyerPasswordResetRepositoryWithReturn) DeleteAllByBuyerID(ctx context.Context, buyerID int) error {
	return nil
}
