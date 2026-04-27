package auction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

type mockAuctionRepository struct {
	auctions []model.Auction
	err      error
}

func (m *mockAuctionRepository) Create(_ context.Context, _ *model.Auction) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepository) FindByID(_ context.Context, _ int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepository) FindByIDWithLock(_ context.Context, _ int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepository) List(_ context.Context, _ *repository.AuctionFilters) ([]model.Auction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.auctions, nil
}
func (m *mockAuctionRepository) ListByVenue(_ context.Context, _ int) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepository) Update(_ context.Context, _ *model.Auction) error {
	return nil
}
func (m *mockAuctionRepository) UpdateStatus(_ context.Context, _ int, _ model.AuctionStatus) error {
	return nil
}
func (m *mockAuctionRepository) Delete(_ context.Context, _ int) error {
	return nil
}

func TestListAuctionsUseCase_Execute(t *testing.T) {
	mockClock := mock.NewMockClock(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))

	auctions := []model.Auction{
		{ID: 1, Status: model.AuctionStatusScheduled, Period: model.NewAuctionPeriod(nil, nil)},
		{ID: 2, Status: model.AuctionStatusInProgress, Period: model.NewAuctionPeriod(nil, nil)},
	}

	tests := []struct {
		name         string
		mockAuctions []model.Auction
		mockErr      error
		wantCount    int
		wantErr      bool
	}{
		{
			name:         "Success",
			mockAuctions: auctions,
			wantCount:    2,
		},
		{
			name:    "RepoError",
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAuctionRepository{auctions: tt.mockAuctions, err: tt.mockErr}
			uc := auction.NewListAuctionsUseCase(repo, mockClock)

			got, err := uc.Execute(context.Background(), nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("got count %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}
