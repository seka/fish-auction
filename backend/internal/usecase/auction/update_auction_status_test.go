package auction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
)

type mockAuctionRepoForStatusUpdate struct {
	err error
}

func (m *mockAuctionRepoForStatusUpdate) Create(ctx context.Context, a *model.Auction) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForStatusUpdate) GetByID(ctx context.Context, id int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForStatusUpdate) List(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForStatusUpdate) ListByVenue(ctx context.Context, venueID int) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForStatusUpdate) Update(ctx context.Context, auction *model.Auction) error {
	return nil
}
func (m *mockAuctionRepoForStatusUpdate) UpdateStatus(ctx context.Context, id int, status model.AuctionStatus) error {
	return m.err
}
func (m *mockAuctionRepoForStatusUpdate) Delete(ctx context.Context, id int) error { return nil }

func TestUpdateAuctionStatusUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		status  model.AuctionStatus
		mockErr error
		wantErr bool
	}{
		{
			name:   "Success",
			id:     1,
			status: model.AuctionStatusInProgress,
		},
		{
			name:    "InvalidStatus",
			id:      1,
			status:  "invalid_status",
			wantErr: true,
		},
		{
			name:    "RepoError",
			id:      1,
			status:  model.AuctionStatusCompleted,
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAuctionRepoForStatusUpdate{err: tt.mockErr}
			uc := auction.NewUpdateAuctionStatusUseCase(repo)

			err := uc.Execute(context.Background(), tt.id, tt.status)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "InvalidStatus" && err != nil {
				expectedMsg := "invalid auction status: " + string(tt.status)
				if err.Error() != expectedMsg {
					t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
				}
			}
		})
	}
}
