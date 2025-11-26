package model_test

import (
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

func TestItemStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status model.ItemStatus
		want   bool
	}{
		{
			name:   "Available",
			status: model.ItemStatusAvailable,
			want:   true,
		},
		{
			name:   "Sold",
			status: model.ItemStatusSold,
			want:   true,
		},
		{
			name:   "Invalid_Empty",
			status: model.ItemStatus(""),
			want:   false,
		},
		{
			name:   "Invalid_Random",
			status: model.ItemStatus("Random"),
			want:   false,
		},
		{
			name:   "Invalid_Lowercase",
			status: model.ItemStatus("available"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status model.ItemStatus
		want   string
	}{
		{
			name:   "Available",
			status: model.ItemStatusAvailable,
			want:   "Available",
		},
		{
			name:   "Sold",
			status: model.ItemStatusSold,
			want:   "Sold",
		},
		{
			name:   "Custom",
			status: model.ItemStatus("Custom"),
			want:   "Custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
