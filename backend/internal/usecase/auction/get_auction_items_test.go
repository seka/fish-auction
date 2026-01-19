package auction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
)

type mockItemRepository struct {
	items []model.AuctionItem
	err   error
}

// Implement ItemRepository interface
func (m *mockItemRepository) Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	return nil, nil
}
func (m *mockItemRepository) FindByID(ctx context.Context, id int) (*model.AuctionItem, error) {
	return nil, nil
}
func (m *mockItemRepository) List(ctx context.Context, status string) ([]model.AuctionItem, error) {
	return nil, nil
}
func (m *mockItemRepository) ListByAuction(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.items, nil
}
func (m *mockItemRepository) Update(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	return nil, nil
}
func (m *mockItemRepository) Delete(ctx context.Context, id int) error {
	return nil
}
func (m *mockItemRepository) UpdateStatus(ctx context.Context, id int, status model.ItemStatus) error {
	return nil
}
func (m *mockItemRepository) UpdateSortOrder(ctx context.Context, id int, sortOrder int) error {
	return nil
}
func (m *mockItemRepository) InvalidateCache(ctx context.Context, id int) error {
	return nil
}

func TestGetAuctionItemsUseCase_Execute(t *testing.T) {
	items := []model.AuctionItem{
		{ID: 1, FishType: "Tuna"},
		{ID: 2, FishType: "Salmon"},
	}

	tests := []struct {
		name      string
		auctionID int
		mockItems []model.AuctionItem
		mockErr   error
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Success",
			auctionID: 1,
			mockItems: items,
			wantCount: 2,
		},
		{
			name:      "RepoError",
			auctionID: 1,
			mockErr:   errors.New("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockItemRepository{items: tt.mockItems, err: tt.mockErr}
			uc := auction.NewGetAuctionItemsUseCase(repo)

			got, err := uc.Execute(context.Background(), tt.auctionID)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("got count %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}
