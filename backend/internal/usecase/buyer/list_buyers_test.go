package buyer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestListBuyersUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		buyers  []model.Buyer
		wantErr error
	}{
		{
			name: "Success",
			buyers: []model.Buyer{
				{ID: 1, Name: "A"},
				{ID: 2, Name: "B"},
			},
		},
		{
			name:    "Error",
			wantErr: errors.New("list failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockBuyerRepository{
				ListFunc: func(ctx context.Context) ([]model.Buyer, error) {
					if tt.wantErr != nil {
						return nil, tt.wantErr
					}
					return tt.buyers, nil
				},
			}

			uc := buyer.NewListBuyersUseCase(repo)
			got, err := uc.Execute(context.Background())

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(got) != len(tt.buyers) {
				t.Fatalf("expected %d buyers, got %d", len(tt.buyers), len(got))
			}
		})
	}
}
