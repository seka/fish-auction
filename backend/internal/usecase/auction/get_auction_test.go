package auction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
)

type mockAuctionRepoForGet struct {
	auction *model.Auction
	err     error
}

func (m *mockAuctionRepoForGet) Create(_ context.Context, _ *model.Auction) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForGet) FindByID(_ context.Context, id int) (*model.Auction, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.auction != nil && m.auction.ID == id {
		return m.auction, nil
	}
	return nil, nil
}
func (m *mockAuctionRepoForGet) FindByIDWithLock(ctx context.Context, id int) (*model.Auction, error) {
	return m.FindByID(ctx, id)
}
func (m *mockAuctionRepoForGet) List(_ context.Context, _ *repository.AuctionFilters) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForGet) ListByVenue(_ context.Context, _ int) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForGet) Update(_ context.Context, _ *model.Auction) error { return nil }
func (m *mockAuctionRepoForGet) UpdateStatus(_ context.Context, _ int, _ model.AuctionStatus) error {
	return nil
}
func (m *mockAuctionRepoForGet) Delete(_ context.Context, _ int) error { return nil }

func TestGetAuctionUseCase_Execute(t *testing.T) {
	validAuction := &model.Auction{ID: 1}

	tests := []struct {
		name    string
		id      int
		mockAuc *model.Auction
		mockErr error
		wantErr bool
		wantNil bool
	}{
		{
			name:    "Success",
			id:      1,
			mockAuc: validAuction,
		},
		{
			name:    "NotFound",
			id:      99,
			mockAuc: nil,
			mockErr: errors.New("not found"),
			wantErr: true,
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
			repo := &mockAuctionRepoForGet{
				auction: tt.mockAuc,
				err:     tt.mockErr,
			}
			uc := auction.NewGetAuctionUseCase(repo)

			got, err := uc.Execute(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNil && got != nil {
				t.Error("expected nil, got auction")
			}
			if !tt.wantNil && !tt.wantErr && got == nil {
				t.Error("expected auction, got nil")
			}
		})
	}
}
