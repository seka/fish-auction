package auction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
)

type mockAuctionRepository struct {
	auctions []model.Auction
	err      error
}

func (m *mockAuctionRepository) Create(ctx context.Context, auction *model.Auction) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepository) GetByID(ctx context.Context, id int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepository) List(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.auctions, nil
}
func (m *mockAuctionRepository) ListByVenue(ctx context.Context, venueID int) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepository) Update(ctx context.Context, auction *model.Auction) error {
	return nil
}
func (m *mockAuctionRepository) UpdateStatus(ctx context.Context, id int, status model.AuctionStatus) error {
	return nil
}
func (m *mockAuctionRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func TestListAuctionsUseCase_Execute(t *testing.T) {
	// Past date setup for comprehensive testing
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	pastDate := time.Now().In(jst).AddDate(0, 0, -1)

	auctions := []model.Auction{
		{ID: 1, Status: model.AuctionStatusScheduled},
		{ID: 2, Status: model.AuctionStatusInProgress, AuctionDate: pastDate}, // Should be completed logic?
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
			uc := auction.NewListAuctionsUseCase(repo)

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
