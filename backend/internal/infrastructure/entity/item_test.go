package entity_test

import (
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

func TestAuctionItem_Validate(t *testing.T) {
	tests := []struct {
		name      string
		item      *entity.AuctionItem
		wantErr   bool
		wantField string
	}{
		{
			name: "Valid",
			item: &entity.AuctionItem{
				FishermanID: 1,
				FishType:    "Tuna",
				Quantity:    10,
				Unit:        "kg",
				Status:      model.ItemStatusAvailable,
			},
		},
		{
			name: "Invalid_FishermanID_Zero",
			item: &entity.AuctionItem{
				FishermanID: 0,
				FishType:    "Tuna",
				Quantity:    10,
				Unit:        "kg",
				Status:      model.ItemStatusAvailable,
			},
			wantErr:   true,
			wantField: "fisherman_id",
		},
		{
			name: "Invalid_FishermanID_Negative",
			item: &entity.AuctionItem{
				FishermanID: -1,
				FishType:    "Tuna",
				Quantity:    10,
				Unit:        "kg",
				Status:      model.ItemStatusAvailable,
			},
			wantErr:   true,
			wantField: "fisherman_id",
		},
		{
			name: "Invalid_FishType_Empty",
			item: &entity.AuctionItem{
				FishermanID: 1,
				FishType:    "",
				Quantity:    10,
				Unit:        "kg",
				Status:      model.ItemStatusAvailable,
			},
			wantErr:   true,
			wantField: "fish_type",
		},
		{
			name: "Invalid_FishType_Whitespace",
			item: &entity.AuctionItem{
				FishermanID: 1,
				FishType:    "   ",
				Quantity:    10,
				Unit:        "kg",
				Status:      model.ItemStatusAvailable,
			},
			wantErr:   true,
			wantField: "fish_type",
		},
		{
			name: "Invalid_Quantity_Zero",
			item: &entity.AuctionItem{
				FishermanID: 1,
				FishType:    "Tuna",
				Quantity:    0,
				Unit:        "kg",
				Status:      model.ItemStatusAvailable,
			},
			wantErr:   true,
			wantField: "quantity",
		},
		{
			name: "Invalid_Unit_Empty",
			item: &entity.AuctionItem{
				FishermanID: 1,
				FishType:    "Tuna",
				Quantity:    10,
				Unit:        "",
				Status:      model.ItemStatusAvailable,
			},
			wantErr:   true,
			wantField: "unit",
		},
		{
			name: "Invalid_Status",
			item: &entity.AuctionItem{
				FishermanID: 1,
				FishType:    "Tuna",
				Quantity:    10,
				Unit:        "kg",
				Status:      model.ItemStatus("Invalid"),
			},
			wantErr:   true,
			wantField: "status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
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

func TestAuctionItem_ToModel(t *testing.T) {
	item := &entity.AuctionItem{
		ID:          1,
		FishermanID: 2,
		FishType:    "Tuna",
		Quantity:    10,
		Unit:        "kg",
		Status:      model.ItemStatusAvailable,
	}

	modelItem := item.ToModel()

	if modelItem.ID != item.ID {
		t.Errorf("expected ID %d, got %d", item.ID, modelItem.ID)
	}
	if modelItem.FishermanID != item.FishermanID {
		t.Errorf("expected FishermanID %d, got %d", item.FishermanID, modelItem.FishermanID)
	}
	if modelItem.FishType != item.FishType {
		t.Errorf("expected FishType %s, got %s", item.FishType, modelItem.FishType)
	}
	if modelItem.Quantity != item.Quantity {
		t.Errorf("expected Quantity %d, got %d", item.Quantity, modelItem.Quantity)
	}
	if modelItem.Unit != item.Unit {
		t.Errorf("expected Unit %s, got %s", item.Unit, modelItem.Unit)
	}
	if modelItem.Status != item.Status {
		t.Errorf("expected Status %s, got %s", item.Status, modelItem.Status)
	}
}
