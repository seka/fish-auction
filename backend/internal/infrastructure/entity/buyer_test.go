package entity_test

import (
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

func TestBuyer_Validate(t *testing.T) {
	tests := []struct {
		name      string
		buyer     *entity.Buyer
		wantErr   bool
		wantField string
	}{
		{
			name: "Valid",
			buyer: &entity.Buyer{
				Name:         "John Doe",
				Organization: "Fish Corp",
				ContactInfo:  "john@example.com",
			},
		},
		{
			name: "Invalid_Name_Empty",
			buyer: &entity.Buyer{
				Name: "",
			},
			wantErr:   true,
			wantField: "name",
		},
		{
			name: "Invalid_Name_Whitespace",
			buyer: &entity.Buyer{
				Name: "   ",
			},
			wantErr:   true,
			wantField: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.buyer.Validate()
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

func TestBuyer_ToModel(t *testing.T) {
	buyer := &entity.Buyer{
		ID:   1,
		Name: "John Doe",
	}

	modelBuyer := buyer.ToModel()

	if modelBuyer.ID != buyer.ID {
		t.Errorf("expected ID %d, got %d", buyer.ID, modelBuyer.ID)
	}
	if modelBuyer.Name != buyer.Name {
		t.Errorf("expected Name %s, got %s", buyer.Name, modelBuyer.Name)
	}
}
