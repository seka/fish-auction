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

type mockAuctionRepoForCreate struct {
	created *model.Auction
	err     error
}

func (m *mockAuctionRepoForCreate) Create(_ context.Context, a *model.Auction) (*model.Auction, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.created = a
	// Simulate ID assignment
	a.ID = 1
	return a, nil
}
func (m *mockAuctionRepoForCreate) FindByID(_ context.Context, _ int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForCreate) FindByIDWithLock(_ context.Context, _ int) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForCreate) List(_ context.Context, _ *repository.AuctionFilters) ([]model.Auction, error) {
	// Actually interface compliance check:
	// repository.AuctionRepository uses repository.AuctionFilters.
	// But minimal mock for Create doesn't need List.
	return nil, nil
}

// Quick minimal methods stub to satisfy interface
func (m *mockAuctionRepoForCreate) ListByVenue(_ context.Context, _ int) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForCreate) Update(_ context.Context, _ *model.Auction) error {
	return nil
}
func (m *mockAuctionRepoForCreate) UpdateStatus(_ context.Context, _ int, _ model.AuctionStatus) error {
	return nil
}
func (m *mockAuctionRepoForCreate) Delete(_ context.Context, _ int) error { return nil }

// Fix List signature
// But we need to import repository package for the signature.
// Alternatively, embed the full mock if we extract it to a shared file, but for now local stub.

func TestCreateAuctionUseCase_Execute(t *testing.T) {
	validAuction := &model.Auction{
		VenueID: 1,
		Status:  model.AuctionStatusScheduled,
		Period:  model.NewAuctionPeriod(time.Now(), nil, nil),
	}

	tests := []struct {
		name    string
		input   *model.Auction
		repoErr error
		wantErr bool
	}{
		{
			name:  "Success",
			input: validAuction,
		},
		{
			name:    "RepoError",
			input:   validAuction,
			repoErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAuctionRepoForCreate{err: tt.repoErr}
			// Note: Assuming NewCreateAuctionUseCase takes repo
			uc := auction.NewCreateAuctionUseCase(repo)

			got, err := uc.Execute(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Error("expected auction, got nil")
			}
			if !tt.wantErr && got.Status != model.AuctionStatusScheduled {
				t.Errorf("expected status scheduled, got %v", got.Status)
			}
		})
	}
}

// Full interface satisfaction needed for the type assertion in constructor?
// Go interfaces are implicit.
// But we need to implement List(ctx, *repository.AuctionFilters)
// So we need to import repository.
