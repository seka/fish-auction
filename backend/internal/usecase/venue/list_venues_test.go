package venue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

type mockVenueRepository struct {
	venues []model.Venue
	err    error
}

func (m *mockVenueRepository) Create(_ context.Context, _ *model.Venue) (*model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepository) FindByID(_ context.Context, _ int) (*model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepository) List(_ context.Context) ([]model.Venue, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.venues, nil
}
func (m *mockVenueRepository) Update(_ context.Context, _ *model.Venue) error { return nil }
func (m *mockVenueRepository) Delete(_ context.Context, _ int) error             { return nil }

func TestListVenuesUseCase_Execute(t *testing.T) {
	venues := []model.Venue{
		{ID: 1, Name: "Venue A"},
		{ID: 2, Name: "Venue B"},
	}

	tests := []struct {
		name       string
		mockVenues []model.Venue
		mockErr    error
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "Success",
			mockVenues: venues,
			wantCount:  2,
		},
		{
			name:    "RepoError",
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockVenueRepository{venues: tt.mockVenues, err: tt.mockErr}
			uc := venue.NewListVenuesUseCase(repo)

			got, err := uc.Execute(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("got count %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}
