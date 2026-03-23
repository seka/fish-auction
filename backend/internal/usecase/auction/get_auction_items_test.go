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
func (m *mockItemRepository) Create(_ context.Context, _ *model.AuctionItem) (*model.AuctionItem, error) {
	return nil, nil
}
func (m *mockItemRepository) FindByID(_ context.Context, _ int) (*model.AuctionItem, error) {
	return nil, nil
}
func (m *mockItemRepository) FindByIDWithLock(_ context.Context, _ int) (*model.AuctionItem, error) {
	return nil, nil
}
func (m *mockItemRepository) List(_ context.Context, _ string) ([]model.AuctionItem, error) {
	return nil, nil
}
func (m *mockItemRepository) ListByAuction(_ context.Context, _ int) ([]model.AuctionItem, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.items, nil
}
func (m *mockItemRepository) Update(_ context.Context, _ *model.AuctionItem) (*model.AuctionItem, error) {
	return nil, nil
}
func (m *mockItemRepository) Delete(_ context.Context, _ int) error {
	return nil
}
func (m *mockItemRepository) UpdateStatus(_ context.Context, _ int, _ model.ItemStatus) error {
	return nil
}
func (m *mockItemRepository) UpdateSortOrder(_ context.Context, _, _ int) error {
	return nil
}
func (m *mockItemRepository) Reorder(_ context.Context, _ int, _ []int) error {
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
