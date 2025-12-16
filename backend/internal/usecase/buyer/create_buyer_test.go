package buyer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
	"golang.org/x/crypto/bcrypt"
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
	}{
		{
			name:     "Success",
			input:    "John Doe",
			password: "password",
			wantID:   10,
		},
		{
			name:           "BuyerRepoError",
			input:          "Jane Doe",
			password:       "password",
			createBuyerErr: createErr,
			wantErr:        createErr,
		},
		{
			name:          "AuthRepoError",
			input:         "John Doe",
			password:      "password",
			wantID:        11,
			createAuthErr: authErr,
			wantErr:       authErr,
		},
		{
			name:     "PasswordTooLong",
			input:    "John Doe",
			password: "this_password_is_definitely_way_too_long_to_be_hashed_by_bcrypt_because_it_exceeds_seventy_two_bytes_limit",
			wantErr:  bcrypt.ErrPasswordTooLong, // Or check functionality logic which probably propagates expected error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buyerRepo := &mock.MockBuyerRepository{
				CreateFunc: func(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error) {
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
				CreateFunc: func(ctx context.Context, auth *model.Authentication) (*model.Authentication, error) {
					if tt.createAuthErr != nil {
						return nil, tt.createAuthErr
					}
					return auth, nil
				},
			}

			uc := buyer.NewCreateBuyerUseCase(buyerRepo, authRepo)
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

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got == nil || got.ID != tt.wantID || got.Name != tt.input {
				t.Fatalf("unexpected buyer %+v", got)
			}
		})
	}
}
