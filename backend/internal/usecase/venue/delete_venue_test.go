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

func (m *mockVenueRepoForDelete) Create(ctx context.Context, v *model.Venue) (*model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepoForDelete) GetByID(ctx context.Context, id int) (*model.Venue, error) {
	return nil, nil
}
func (m *mockVenueRepoForDelete) List(ctx context.Context) ([]model.Venue, error) { return nil, nil }
func (m *mockVenueRepoForDelete) Update(ctx context.Context, venue *model.Venue) error {
	return nil
}
func (m *mockVenueRepoForDelete) Delete(ctx context.Context, id int) error {
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
