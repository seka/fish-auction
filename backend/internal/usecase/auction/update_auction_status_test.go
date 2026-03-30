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

func (m *mockAuctionRepoForStatusUpdate) Create(_ context.Context, _ *model.Auction) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForStatusUpdate) FindByID(_ context.Context, _ int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForStatusUpdate) FindByIDWithLock(_ context.Context, _ int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForStatusUpdate) List(_ context.Context, _ *repository.AuctionFilters) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForStatusUpdate) ListByVenue(_ context.Context, _ int) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForStatusUpdate) Update(_ context.Context, _ *model.Auction) error {
	return nil
}
func (m *mockAuctionRepoForStatusUpdate) UpdateStatus(_ context.Context, _ int, _ model.AuctionStatus) error {
	return m.err
}
func (m *mockAuctionRepoForStatusUpdate) Delete(_ context.Context, _ int) error { return nil }

type mockBuyerRepoForStatusUpdate struct{}

func (m *mockBuyerRepoForStatusUpdate) Create(_ context.Context, _ *model.Buyer) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockBuyerRepoForStatusUpdate) List(_ context.Context) ([]model.Buyer, error) {
	return []model.Buyer{{ID: 1}}, nil
}
func (m *mockBuyerRepoForStatusUpdate) FindByID(_ context.Context, _ int) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockBuyerRepoForStatusUpdate) FindByName(_ context.Context, _ string) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockBuyerRepoForStatusUpdate) FindByEmail(_ context.Context, _ string) (*model.Buyer, error) {
	return nil, nil
}
func (m *mockBuyerRepoForStatusUpdate) Delete(_ context.Context, _ int) error { return nil }

type mockPushUseCaseForStatusUpdate struct{}

func (m *mockPushUseCaseForStatusUpdate) Subscribe(_ context.Context, _ int, _ *model.PushSubscription) error {
	return nil
}
func (m *mockPushUseCaseForStatusUpdate) SendNotification(_ context.Context, _ int, _ any) error {
	return nil
}

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
			buyerRepo := &mockBuyerRepoForStatusUpdate{}
			pushUseCase := &mockPushUseCaseForStatusUpdate{}
			uc := auction.NewUpdateAuctionStatusUseCase(repo, buyerRepo, pushUseCase)

			err := uc.Execute(context.Background(), tt.id, tt.status)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "InvalidStatus" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var invalidStatusErr *auction.InvalidStatusError
				if !errors.As(err, &invalidStatusErr) {
					t.Fatalf("expected InvalidStatusError, got %T: %v", err, err)
				}
				if invalidStatusErr.Status != string(tt.status) {
					t.Errorf("expected status %q in error, got %q", string(tt.status), invalidStatusErr.Status)
				}
			}
		})
	}
}
