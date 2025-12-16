package buyer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginBuyerUseCase_Execute(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	validAuth := &model.Authentication{
		ID:           1,
		BuyerID:      1,
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}
	validBuyer := &model.Buyer{ID: 1}

	lockedUntil := time.Now().Add(1 * time.Hour)
	lockedAuth := &model.Authentication{
		ID:           2,
		BuyerID:      2,
		Email:        "locked@example.com",
		PasswordHash: string(hashedPassword),
		LockedUntil:  &lockedUntil,
	}

	maxFailedAuth := &model.Authentication{
		ID:             3,
		BuyerID:        3,
		Email:          "fail@example.com",
		PasswordHash:   string(hashedPassword),
		FailedAttempts: 4, // Next fail will lock
	}

	tests := []struct {
		name           string
		email          string
		password       string
		mockAuth       *model.Authentication
		mockBuyer      *model.Buyer
		mockAuthErr    error
		mockBuyerErr   error
		wantErr        bool
		wantLockCalled bool
		wantLocked     bool
	}{
		{
			name:      "Success",
			email:     "test@example.com",
			password:  "password",
			mockAuth:  validAuth,
			mockBuyer: validBuyer,
		},
		{
			name:        "InvalidEmail",
			email:       "wrong@example.com",
			password:    "password",
			mockAuth:    nil,
			mockAuthErr: errors.New("not found"), // Repo returns error on not found usually? Or nil?
			// Usecase: auth, err := uc.authRepo.FindByEmail(ctx, email); if err != nil ...
			// If implementation treats err as invalid credentials, then we simulate err.
			wantErr: true,
		},
		{
			name:     "InvalidPassword",
			email:    "test@example.com",
			password: "wrong",
			mockAuth: validAuth,
			wantErr:  true,
		},
		{
			name:     "AccountLocked",
			email:    "locked@example.com",
			password: "password", // Password might be correct but account is locked
			mockAuth: lockedAuth,
			wantErr:  true,
		},
		{
			name:           "LockoutTriggered",
			email:          "fail@example.com",
			password:       "wrong",
			mockAuth:       maxFailedAuth,
			wantErr:        true,
			wantLockCalled: true,
		},
		{
			name:         "RepoError_FindByID",
			email:        "test@example.com",
			password:     "password",
			mockAuth:     validAuth,
			mockBuyerErr: errors.New("db error"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockCalled := false

			mockAuthRepo := &mock.MockAuthenticationRepository{
				FindByEmailFunc: func(ctx context.Context, email string) (*model.Authentication, error) {
					if tt.mockAuthErr != nil {
						return nil, tt.mockAuthErr
					}
					if tt.mockAuth != nil && tt.mockAuth.Email == email {
						return tt.mockAuth, nil
					}
					return nil, errors.New("not found")
				},
				UpdateLoginSuccessFunc: func(ctx context.Context, id int, loginAt time.Time) error {
					return nil
				},
				IncrementFailedAttemptsFunc: func(ctx context.Context, id int) error {
					return nil
				},
				LockAccountFunc: func(ctx context.Context, id int, until time.Time) error {
					lockCalled = true
					return nil
				},
			}

			mockBuyerRepo := &mock.MockBuyerRepository{
				FindByIDFunc: func(ctx context.Context, id int) (*model.Buyer, error) {
					if tt.mockBuyerErr != nil {
						return nil, tt.mockBuyerErr
					}
					return tt.mockBuyer, nil
				},
			}

			uc := buyer.NewLoginBuyerUseCase(mockBuyerRepo, mockAuthRepo)
			_, err := uc.Execute(context.Background(), tt.email, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantLockCalled != lockCalled {
				t.Errorf("lockCalled = %v, want %v", lockCalled, tt.wantLockCalled)
			}
		})
	}
}
