package buyer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestCreateBuyerUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantID  int
		wantErr error
	}{
		{
			name:   "Success",
			input:  "John Doe",
			wantID: 10,
		},
		{
			name:    "Error",
			input:   "Jane Doe",
			wantErr: errors.New("create failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buyerRepo := &mock.MockBuyerRepository{
				CreateFunc: func(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error) {
					if buyer.Name != tt.input {
						t.Fatalf("unexpected name %s", buyer.Name)
					}
					if tt.wantErr != nil {
						return nil, tt.wantErr
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
					if tt.wantErr != nil {
						return nil, tt.wantErr
					}
					return auth, nil
				},
			}

			uc := buyer.NewCreateBuyerUseCase(buyerRepo, authRepo)
			got, err := uc.Execute(context.Background(), tt.input, "test@example.com", "password", "org", "contact")

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
