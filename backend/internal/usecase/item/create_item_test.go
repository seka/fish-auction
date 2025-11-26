package item_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestCreateItemUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		input   *model.AuctionItem
		wantID  int
		wantErr error
	}{
		{
			name: "Success",
			input: &model.AuctionItem{
				FishermanID: 1,
				FishType:    "Tuna",
				Quantity:    10,
				Unit:        "kg",
				Status:      model.ItemStatusAvailable,
			},
			wantID: 1,
		},
		{
			name: "Error",
			input: &model.AuctionItem{
				FishermanID: 1,
				FishType:    "Tuna",
				Quantity:    10,
				Unit:        "kg",
				Status:      model.ItemStatusAvailable,
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockItemRepository{
				CreateFunc: func(ctx context.Context, in *model.AuctionItem) (*model.AuctionItem, error) {
					if in != tt.input {
						t.Fatalf("input pointer mismatch")
					}
					if tt.wantErr != nil {
						return nil, tt.wantErr
					}
					cloned := *in
					cloned.ID = tt.wantID
					cloned.CreatedAt = time.Now()
					return &cloned, nil
				},
			}

			uc := item.NewCreateItemUseCase(mockRepo)
			created, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if created != nil {
					t.Fatalf("expected nil result, got %+v", created)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if created == nil || created.ID != tt.wantID {
				t.Fatalf("unexpected created item %+v", created)
			}
		})
	}
}
