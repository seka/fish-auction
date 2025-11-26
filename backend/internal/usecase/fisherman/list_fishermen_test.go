package fisherman_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/fisherman"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestListFishermenUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		fishermen []model.Fisherman
		wantErr   error
	}{
		{
			name: "Success",
			fishermen: []model.Fisherman{
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
			repo := &mock.MockFishermanRepository{
				ListFunc: func(ctx context.Context) ([]model.Fisherman, error) {
					if tt.wantErr != nil {
						return nil, tt.wantErr
					}
					return tt.fishermen, nil
				},
			}

			uc := fisherman.NewListFishermenUseCase(repo)
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
			if len(got) != len(tt.fishermen) {
				t.Fatalf("expected %d fishermen, got %d", len(tt.fishermen), len(got))
			}
		})
	}
}
