package venue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

type mockVenueRepoForDelete struct {
	err error
}

func (m *mockVenueRepoForDelete) Create(_ context.Context, _ *model.Venue) (*model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepoForDelete) FindByID(_ context.Context, _ int) (*model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepoForDelete) List(_ context.Context) ([]model.Venue, error) { return nil, nil }
func (m *mockVenueRepoForDelete) Update(_ context.Context, _ *model.Venue) error {
	return nil
}
func (m *mockVenueRepoForDelete) Delete(_ context.Context, _ int) error {
	return m.err
}

func TestDeleteVenueUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		repoErr error
		wantErr bool
	}{
		{
			name: "Success",
			id:   1,
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
			repo := &mockVenueRepoForDelete{err: tt.repoErr}
			uc := venue.NewDeleteVenueUseCase(repo)

			err := uc.Execute(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
