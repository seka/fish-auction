package venue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

type mockVenueRepoForGet struct {
	venue *model.Venue
	err   error
}

func (m *mockVenueRepoForGet) Create(ctx context.Context, v *model.Venue) (*model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepoForGet) GetByID(ctx context.Context, id int) (*model.Venue, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.venue != nil && m.venue.ID == id {
		return m.venue, nil
	}
	return nil, nil
}
func (m *mockVenueRepoForGet) List(ctx context.Context) ([]model.Venue, error)      { return nil, nil }
func (m *mockVenueRepoForGet) Update(ctx context.Context, venue *model.Venue) error { return nil }
func (m *mockVenueRepoForGet) Delete(ctx context.Context, id int) error             { return nil }

func TestGetVenueUseCase_Execute(t *testing.T) {
	validVenue := &model.Venue{ID: 1, Name: "Venue A"}

	tests := []struct {
		name      string
		id        int
		mockVenue *model.Venue
		repoErr   error
		wantErr   bool
		wantNil   bool
	}{
		{
			name:      "Success",
			id:        1,
			mockVenue: validVenue,
		},
		{
			name:      "NotFound",
			id:        99,
			mockVenue: validVenue,
			wantNil:   true,
		},
		{
			name:    "RepoError",
			id:      1,
			repoErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockVenueRepoForGet{venue: tt.mockVenue, err: tt.repoErr}
			uc := venue.NewGetVenueUseCase(repo)

			got, err := uc.Execute(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNil && got != nil {
				t.Error("expected nil, got venue")
			}
			if !tt.wantNil && !tt.wantErr && got == nil {
				t.Error("expected venue, got nil")
			}
		})
	}
}
