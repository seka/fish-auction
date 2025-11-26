package entity_test

import (
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

func TestFisherman_Validate(t *testing.T) {
	tests := []struct {
		name      string
		fisherman *entity.Fisherman
		wantErr   bool
		wantField string
	}{
		{
			name: "Valid",
			fisherman: &entity.Fisherman{
				Name: "Captain Jack",
			},
		},
		{
			name: "Invalid_Name_Empty",
			fisherman: &entity.Fisherman{
				Name: "",
			},
			wantErr:   true,
			wantField: "name",
		},
		{
			name: "Invalid_Name_Whitespace",
			fisherman: &entity.Fisherman{
				Name: "   ",
			},
			wantErr:   true,
			wantField: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fisherman.Validate()
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

func TestFisherman_ToModel(t *testing.T) {
	fisherman := &entity.Fisherman{
		ID:   1,
		Name: "Captain Jack",
	}

	modelFisherman := fisherman.ToModel()

	if modelFisherman.ID != fisherman.ID {
		t.Errorf("expected ID %d, got %d", fisherman.ID, modelFisherman.ID)
	}
	if modelFisherman.Name != fisherman.Name {
		t.Errorf("expected Name %s, got %s", fisherman.Name, modelFisherman.Name)
	}
}
