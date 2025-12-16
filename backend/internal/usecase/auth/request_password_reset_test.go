package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
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
func (m *mockBuyerRepository) Count(ctx context.Context) (int, error) { return 0, nil }

type mockBuyerPasswordResetRepository struct {
	createErr error
	deleteErr error
}

func (m *mockBuyerPasswordResetRepository) Create(ctx context.Context, buyerID int, tokenHash string, expiresAt time.Time) error {
	return m.createErr
}
func (m *mockBuyerPasswordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (int, time.Time, error) {
	return 0, time.Time{}, nil
}
func (m *mockBuyerPasswordResetRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	return nil
}
func (m *mockBuyerPasswordResetRepository) DeleteAllByBuyerID(ctx context.Context, buyerID int) error {
	return m.deleteErr
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
			name:      "DeleteTokenError",
			email:     "buyer@example.com",
			mockBuyer: validBuyer,
			mockResetRepo: &mockBuyerPasswordResetRepository{
				deleteErr: errors.New("delete failed"),
			},
			wantError: true,
		},
		{
			name:      "CreateTokenError",
			email:     "buyer@example.com",
			mockBuyer: validBuyer,
			mockResetRepo: &mockBuyerPasswordResetRepository{
				createErr: errors.New("create failed"),
			},
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
			resetRepo := tt.mockResetRepo
			if resetRepo == nil {
				resetRepo = &mockBuyerPasswordResetRepository{}
			}
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
