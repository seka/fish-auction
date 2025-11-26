package fisherman_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/fisherman"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestCreateFishermanUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantID  int
		wantErr error
	}{
		{
			name:   "Success",
			input:  "Captain Jack",
			wantID: 7,
		},
		{
			name:    "Error",
			input:   "Captain Nemo",
			wantErr: errors.New("create failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockFishermanRepository{
				CreateFunc: func(ctx context.Context, name string) (*model.Fisherman, error) {
					if name != tt.input {
						t.Fatalf("unexpected name %s", name)
					}
					if tt.wantErr != nil {
						return nil, tt.wantErr
					}
					return &model.Fisherman{
						ID:   tt.wantID,
						Name: name,
					}, nil
				},
			}

			uc := fisherman.NewCreateFishermanUseCase(repo)
			got, err := uc.Execute(context.Background(), tt.input)

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
				t.Fatalf("unexpected fisherman %+v", got)
			}
		})
	}
}
