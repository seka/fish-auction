package venue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

type mockVenueRepoForCreate struct {
	created *model.Venue
	err     error
}

func (m *mockVenueRepoForCreate) Create(ctx context.Context, v *model.Venue) (*model.Venue, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.created = v
	v.ID = 1
	return v, nil
}
func (m *mockVenueRepoForCreate) GetByID(ctx context.Context, id int) (*model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepoForCreate) List(ctx context.Context) ([]model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepoForCreate) Update(ctx context.Context, venue *model.Venue) error { return nil }
func (m *mockVenueRepoForCreate) Delete(ctx context.Context, id int) error             { return nil }

func TestCreateVenueUseCase_Execute(t *testing.T) {
	validVenue := &model.Venue{Name: "New Venue"}

	tests := []struct {
		name    string
		input   *model.Venue
		repoErr error
		wantErr bool
	}{
		{
			name:  "Success",
			input: validVenue,
		},
		{
			name:    "RepoError",
			input:   validVenue,
			repoErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockVenueRepoForCreate{err: tt.repoErr}
			uc := venue.NewCreateVenueUseCase(repo)

			got, err := uc.Execute(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.ID != 1 {
				t.Errorf("expected ID 1, got %d", got.ID)
			}
		})
	}
}
