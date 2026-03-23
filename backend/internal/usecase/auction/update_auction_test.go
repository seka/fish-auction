package auction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
)

type mockAuctionRepoForUpdate struct {
	err error
}

func (m *mockAuctionRepoForUpdate) Create(_ context.Context, _ *model.Auction) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForUpdate) FindByID(_ context.Context, _ int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForUpdate) FindByIDWithLock(_ context.Context, _ int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForUpdate) List(_ context.Context, _ *repository.AuctionFilters) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForUpdate) ListByVenue(_ context.Context, _ int) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForUpdate) Update(_ context.Context, _ *model.Auction) error {
	return m.err
}
func (m *mockAuctionRepoForUpdate) UpdateStatus(_ context.Context, _ int, _ model.AuctionStatus) error {
	return nil
}
func (m *mockAuctionRepoForUpdate) Delete(_ context.Context, _ int) error { return nil }

func TestUpdateAuctionUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		input   *model.Auction
		mockErr error
		wantErr bool
	}{
		{
			name:  "Success",
			input: &model.Auction{ID: 1},
		},
		{
			name:    "RepoError",
			input:   &model.Auction{ID: 1},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAuctionRepoForUpdate{err: tt.mockErr}
			uc := auction.NewUpdateAuctionUseCase(repo)

			err := uc.Execute(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
