package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
	"github.com/stretchr/testify/mock"
)

type mockBuyerRepository struct {
	buyer *model.Buyer
	err   error
}

func (m *mockBuyerRepository) Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockBuyerRepository) FindByID(ctx context.Context, id int) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockBuyerRepository) FindByEmail(ctx context.Context, email string) (*model.Buyer, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Simulate checking email (mapped to buyer in mock setup or logic)
	// Since model.Buyer doesn't have Email, use external check or standard return
	return m.buyer, nil
}
func (m *mockBuyerRepository) List(ctx context.Context) ([]model.Buyer, error) { return nil, nil }
func (m *mockBuyerRepository) FindByName(ctx context.Context, name string) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockBuyerRepository) Count(ctx context.Context) (int, error)   { return 0, nil }
func (m *mockBuyerRepository) Delete(ctx context.Context, id int) error { return nil }

type mockBuyerPasswordResetRepository struct {
	mock.Mock
}

func (m *mockBuyerPasswordResetRepository) Create(ctx context.Context, userID int, role string, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, role, tokenHash, expiresAt)
	return args.Error(0)
}
func (m *mockBuyerPasswordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (int, string, time.Time, error) {
	args := m.Called(ctx, tokenHash)
	return args.Int(0), args.String(1), args.Get(2).(time.Time), args.Error(3)
}
func (m *mockBuyerPasswordResetRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}
func (m *mockBuyerPasswordResetRepository) DeleteAllByUserID(ctx context.Context, userID int, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

type mockEmailService struct {
	sentBuyerURL string
	sentAdminURL string
	err          error
}

func (m *mockEmailService) SendBuyerPasswordReset(ctx context.Context, to, url string) error {
	if m.err != nil {
		return m.err
	}
	m.sentBuyerURL = url
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
	validBuyer := &model.Buyer{ID: 1, Name: "Test Buyer"}

	tests := []struct {
		name          string
		email         string
		mockBuyer     *model.Buyer
		mockRepoErr   error
		mockEmailErr  error
		mockResetRepo *mockBuyerPasswordResetRepository
		wantError     bool
		wantSent      bool
	}{
		{
			name:      "Success",
			email:     "buyer@example.com",
			mockBuyer: validBuyer,
			wantSent:  true,
		},
		{
			name:      "UserNotFound",
			email:     "other@example.com",
			mockBuyer: nil,
			wantSent:  false,
		},
		{
			name:        "RepoError",
			email:       "buyer@example.com",
			mockRepoErr: errors.New("db error"),
			wantError:   false, // Returns nil for security
		},
		{
			name:         "EmailError",
			email:        "buyer@example.com",
			mockBuyer:    validBuyer,
			mockEmailErr: errors.New("email failed"),
			wantError:    true,
		},
		{
			name:          "DeleteTokenError",
			email:         "buyer@example.com",
			mockBuyer:     validBuyer,
			mockResetRepo: nil, // Setup in body
			wantError:     true,
		},
		{
			name:          "CreateTokenError",
			email:         "buyer@example.com",
			mockBuyer:     validBuyer,
			mockResetRepo: nil, // Will setup in test body
			// we can't put mock expectations in struct easily here without changing struct type.
			// Let's rely on name/logic to setup mocks.
			wantError: true,
		},
		{
			name:      "RandomError",
			email:     "buyer@example.com",
			mockBuyer: validBuyer,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Need to add methods to satisfy interface if changed.
			// Re-check interface:
			// List() ([]model.Buyer, error)
			// FindByName(ctx, name) (*model.Buyer, error)
			// Add FindByName to mock if missing?
			// Yes, existing mock is manual and might be missing FindByName which is in interface.

			buyerRepo := &mockBuyerRepository{buyer: tt.mockBuyer, err: tt.mockRepoErr}
			resetRepo := &mockBuyerPasswordResetRepository{}
			if tt.name == "Success" {
				resetRepo.On("DeleteAllByUserID", mock.Anything, 1, "buyer").Return(nil)
				resetRepo.On("Create", mock.Anything, 1, "buyer", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
			} else if tt.name == "DeleteTokenError" {
				resetRepo.On("DeleteAllByUserID", mock.Anything, 1, "buyer").Return(errors.New("delete failed"))
			} else if tt.name == "CreateTokenError" {
				resetRepo.On("DeleteAllByUserID", mock.Anything, 1, "buyer").Return(nil)
				resetRepo.On("Create", mock.Anything, 1, "buyer", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(errors.New("create failed"))
			} else if tt.name == "EmailError" {
				resetRepo.On("DeleteAllByUserID", mock.Anything, 1, "buyer").Return(nil)
				resetRepo.On("Create", mock.Anything, 1, "buyer", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
			}
			// Other cases like UserNotFound don't call repo methods.

			// For RandomError, if it fails before repo, no calls. SetRandRead fails inside Execute?
			// Random error happens during token generation which is before repo calls. So no repo calls expected.
			emailService := &mockEmailService{err: tt.mockEmailErr}

			// Mock rand.Read if testing random error
			if tt.name == "RandomError" {
				cleanup := auth.SetRandRead(func(b []byte) (int, error) {
					return 0, errors.New("random failed")
				})
				defer cleanup()
			}

			uc := auth.NewRequestPasswordResetUseCase(buyerRepo, resetRepo, emailService)
			err := uc.Execute(context.Background(), tt.email)

			if (err != nil) != tt.wantError {
				// Re-check logic: FindByEmail error returns nil.
				if tt.name == "RepoError" && err == nil {
					// Expected behavior as per implementation
				} else {
					t.Errorf("expected error=%v, got %v", tt.wantError, err)
				}
			}

			if tt.wantSent && emailService.sentBuyerURL == "" {
				t.Error("expected email to be sent, but wasn't")
			}
			if !tt.wantSent && emailService.sentBuyerURL != "" {
				t.Error("expected email NOT to be sent, but was")
			}
		})
	}
}
