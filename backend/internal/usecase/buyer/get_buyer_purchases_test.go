package buyer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
)

type mockBidRepoForPurchases struct {
	purchases []model.Purchase
	err       error
}

func (m *mockBidRepoForPurchases) Create(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	return nil, nil
}
func (m *mockBidRepoForPurchases) ListInvoices(ctx context.Context) ([]model.InvoiceItem, error) {
	return nil, nil
}
func (m *mockBidRepoForPurchases) ListByAuctionID(ctx context.Context, auctionID int) ([]model.Bid, error) {
	return nil, nil
}
func (m *mockBidRepoForPurchases) ListAuctionsByBuyerID(ctx context.Context, buyerID int) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockBidRepoForPurchases) ListPurchasesByBuyerID(ctx context.Context, buyerID int) ([]model.Purchase, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.purchases, nil
}
func (m *mockBidRepoForPurchases) GetHighestBid(ctx context.Context, auctionItemID int) (*model.Bid, error) {
	return nil, nil
}

func TestGetBuyerPurchasesUseCase_Execute(t *testing.T) {
	purchases := []model.Purchase{
		{ID: 1, BuyerID: 1, Price: 1000},
		{ID: 2, BuyerID: 1, Price: 2000},
	}

	tests := []struct {
		name      string
		buyerID   int
		mockPurch []model.Purchase
		mockErr   error
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Success",
			buyerID:   1,
			mockPurch: purchases,
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
			repo := &mockBidRepoForPurchases{purchases: tt.mockPurch, err: tt.mockErr}
			uc := buyer.NewGetBuyerPurchasesUseCase(repo)

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
