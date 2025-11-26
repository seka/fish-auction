package item_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestListItemsUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		items   []model.AuctionItem
		wantErr error
	}{
		{
			name:   "Success",
			status: "Available",
			items: []model.AuctionItem{
				{ID: 1, FishType: "Tuna"},
				{ID: 2, FishType: "Salmon"},
			},
		},
		{
			name:    "Error",
			status:  "Sold",
			wantErr: errors.New("list failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockItemRepository{
				ListFunc: func(ctx context.Context, status string) ([]model.AuctionItem, error) {
					if status != tt.status {
						t.Fatalf("unexpected status %s", status)
					}
					if tt.wantErr != nil {
						return nil, tt.wantErr
					}
					return tt.items, nil
				},
			}

			uc := item.NewListItemsUseCase(repo)
			got, err := uc.Execute(context.Background(), tt.status)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(got) != len(tt.items) {
				t.Fatalf("expected %d items, got %d", len(tt.items), len(got))
			}
		})
	}
}
