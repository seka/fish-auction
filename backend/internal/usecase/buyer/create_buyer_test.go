package buyer_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	domainerrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestCreateBuyerUseCase_Execute(t *testing.T) {
	createErr := errors.New("create failed")
	authErr := errors.New("auth failed")

	tests := []struct {
		name           string
		input          string
		password       string
		wantID         int
		createBuyerErr error
		createAuthErr  error
		wantErr        error
		errContains    string
		wantValErr     bool
	}{
		{
			name:     "Success",
			input:    "John Doe",
			password: "Password123",
			wantID:   10,
		},
		{
			name:           "BuyerRepoError",
			input:          "Jane Doe",
			password:       "Password123",
			createBuyerErr: createErr,
			wantErr:        createErr,
		},
		{
			name:          "AuthRepoError",
			input:         "John Doe",
			password:      "Password123",
			wantID:        11,
			createAuthErr: authErr,
			wantErr:       authErr,
		},
		{
			name:        "PasswordTooLong",
			input:       "John Doe",
			password:    "this_password_is_definitely_way_too_long_to_be_hashed_by_bcrypt_because_it_exceeds_seventy_two_bytes_limit",
			wantValErr:  true,
			errContains: "between 8 and 72 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buyerRepo := &mock.MockBuyerRepository{
				CreateFunc: func(_ context.Context, buyer *model.Buyer) (*model.Buyer, error) {
					if buyer.Name != tt.input {
						t.Fatalf("unexpected name %s", buyer.Name)
					}
					if tt.createBuyerErr != nil {
						return nil, tt.createBuyerErr
					}
					return &model.Buyer{
						ID:           tt.wantID,
						Name:         buyer.Name,
						Organization: buyer.Organization,
						ContactInfo:  buyer.ContactInfo,
					}, nil
				},
			}

			authRepo := &mock.MockAuthenticationRepository{
				CreateFunc: func(_ context.Context, auth *model.Authentication) (*model.Authentication, error) {
					if tt.createAuthErr != nil {
						return nil, tt.createAuthErr
					}
					return auth, nil
				},
			}

			txMgr := &mock.MockTransactionManager{}
			uc := buyer.NewCreateBuyerUseCase(buyerRepo, authRepo, txMgr)
			got, err := uc.Execute(context.Background(), tt.input, "test@example.com", tt.password, "org", "contact")

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if got != nil {
					t.Fatalf("expected nil result, got %+v", got)
				}
				return
			}

			if tt.errContains != "" || tt.wantValErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.wantValErr {
					var valErr *domainerrors.ValidationError
					if !errors.As(err, &valErr) {
						t.Fatalf("expected ValidationError, got %T: %v", err, err)
					}
					if valErr.Field != "password" {
						t.Errorf("expected validation error on field 'password', got %q", valErr.Field)
					}
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error containing %q, got %v", tt.errContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got == nil || got.ID != tt.wantID || got.Name != tt.input {
				t.Fatalf("unexpected buyer %+v", got)
			}
		})
	}
}
