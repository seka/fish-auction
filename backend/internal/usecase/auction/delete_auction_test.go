package auction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
)

type mockAuctionRepoForDelete struct {
	err error
}

func (m *mockAuctionRepoForDelete) Create(_ context.Context, _ *model.Auction) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForDelete) FindByID(_ context.Context, _ int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForDelete) FindByIDWithLock(_ context.Context, _ int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForDelete) List(_ context.Context, _ *repository.AuctionFilters) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForDelete) ListByVenue(_ context.Context, _ int) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForDelete) Update(_ context.Context, _ *model.Auction) error {
	return nil
}
func (m *mockAuctionRepoForDelete) UpdateStatus(_ context.Context, _ int, _ model.AuctionStatus) error {
	return nil
}
func (m *mockAuctionRepoForDelete) Delete(_ context.Context, _ int) error {
	return m.err
}

func TestDeleteAuctionUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		mockErr error
		wantErr bool
	}{
		{
			name: "Success",
			id:   1,
		},
		{
			name:    "RepoError",
			id:      1,
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAuctionRepoForDelete{err: tt.mockErr}
			uc := auction.NewDeleteAuctionUseCase(repo)

			err := uc.Execute(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
