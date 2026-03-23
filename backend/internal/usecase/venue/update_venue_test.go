package venue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

type mockVenueRepoForUpdate struct {
	err error
}

func (m *mockVenueRepoForUpdate) Create(_ context.Context, _ *model.Venue) (*model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepoForUpdate) FindByID(_ context.Context, _ int) (*model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepoForUpdate) List(_ context.Context) ([]model.Venue, error) { return nil, nil }
func (m *mockVenueRepoForUpdate) Update(_ context.Context, _ *model.Venue) error {
	return m.err
}
func (m *mockVenueRepoForUpdate) Delete(_ context.Context, _ int) error { return nil }

func TestUpdateVenueUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		input   *model.Venue
		repoErr error
		wantErr bool
	}{
		{
			name:  "Success",
			input: &model.Venue{ID: 1, Name: "Updated"},
		},
		{
			name:    "RepoError",
			input:   &model.Venue{ID: 1},
			repoErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockVenueRepoForUpdate{err: tt.repoErr}
			uc := venue.NewUpdateVenueUseCase(repo)

			err := uc.Execute(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
