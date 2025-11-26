package entity_test

import (
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

func TestBid_Validate(t *testing.T) {
	tests := []struct {
		name      string
		bid       *entity.Bid
		wantErr   bool
		wantField string
	}{
		{
			name: "Valid",
			bid: &entity.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
		},
		{
			name: "Invalid_ItemID_Zero",
			bid: &entity.Bid{
				ItemID:  0,
				BuyerID: 1,
				Price:   1000,
			},
			wantErr:   true,
			wantField: "item_id",
		},
		{
			name: "Invalid_ItemID_Negative",
			bid: &entity.Bid{
				ItemID:  -1,
				BuyerID: 1,
				Price:   1000,
			},
			wantErr:   true,
			wantField: "item_id",
		},
		{
			name: "Invalid_BuyerID_Zero",
			bid: &entity.Bid{
				ItemID:  1,
				BuyerID: 0,
				Price:   1000,
			},
			wantErr:   true,
			wantField: "buyer_id",
		},
		{
			name: "Invalid_Price_Zero",
			bid: &entity.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   0,
			},
			wantErr:   true,
			wantField: "price",
		},
		{
			name: "Invalid_Price_Negative",
			bid: &entity.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   -100,
			},
			wantErr:   true,
			wantField: "price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bid.Validate()
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("expected validation error, got nil")
			}

			validationErr, ok := err.(*errors.ValidationError)
			if !ok {
				t.Fatalf("expected ValidationError, got %T", err)
			}
			if tt.wantField != "" && validationErr.Field != tt.wantField {
				t.Fatalf("expected field %s, got %s", tt.wantField, validationErr.Field)
			}
		})
	}
}

func TestBid_ToModel(t *testing.T) {
	bid := &entity.Bid{
		ID:      1,
		ItemID:  2,
		BuyerID: 3,
		Price:   1000,
	}

	modelBid := bid.ToModel()

	if modelBid.ID != bid.ID {
		t.Errorf("expected ID %d, got %d", bid.ID, modelBid.ID)
	}
	if modelBid.ItemID != bid.ItemID {
		t.Errorf("expected ItemID %d, got %d", bid.ItemID, modelBid.ItemID)
	}
	if modelBid.BuyerID != bid.BuyerID {
		t.Errorf("expected BuyerID %d, got %d", bid.BuyerID, modelBid.BuyerID)
	}
	if modelBid.Price != bid.Price {
		t.Errorf("expected Price %d, got %d", bid.Price, modelBid.Price)
	}
}
