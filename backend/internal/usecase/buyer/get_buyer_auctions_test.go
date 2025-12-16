package buyer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
)

type mockBidRepoForAuctions struct {
	auctions []model.Auction
	err      error
}

func (m *mockBidRepoForAuctions) Create(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	return nil, nil
}
func (m *mockBidRepoForAuctions) ListInvoices(ctx context.Context) ([]model.InvoiceItem, error) {
	return nil, nil
}
func (m *mockBidRepoForAuctions) ListByAuctionID(ctx context.Context, auctionID int) ([]model.Bid, error) {
	return nil, nil
}
func (m *mockBidRepoForAuctions) ListAuctionsByBuyerID(ctx context.Context, buyerID int) ([]model.Auction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.auctions, nil
}
func (m *mockBidRepoForAuctions) ListPurchasesByBuyerID(ctx context.Context, buyerID int) ([]model.Purchase, error) {
	return nil, nil
}
func (m *mockBidRepoForAuctions) GetHighestBid(ctx context.Context, auctionItemID int) (*model.Bid, error) {
	return nil, nil
}

func TestGetBuyerAuctionsUseCase_Execute(t *testing.T) {
	auctions := []model.Auction{
		{ID: 1, Status: model.AuctionStatusScheduled},
		{ID: 2, Status: model.AuctionStatusInProgress},
	}

	tests := []struct {
		name      string
		buyerID   int
		mockAucs  []model.Auction
		mockErr   error
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Success",
			buyerID:   1,
			mockAucs:  auctions,
			wantCount: 2,
		},
		{
			name:    "RepoError",
			buyerID: 1,
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBidRepoForAuctions{auctions: tt.mockAucs, err: tt.mockErr}
			uc := buyer.NewGetBuyerAuctionsUseCase(repo)

			got, err := uc.Execute(context.Background(), tt.buyerID)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("got count %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}
