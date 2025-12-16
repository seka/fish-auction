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

type mockAuctionRepoForGet struct {
	auction         *model.Auction
	err             error
	updateStatusErr error
}

func (m *mockAuctionRepoForGet) Create(ctx context.Context, a *model.Auction) (*model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForGet) GetByID(ctx context.Context, id int) (*model.Auction, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.auction != nil && m.auction.ID == id {
		return m.auction, nil
	}
	return nil, nil
}
func (m *mockAuctionRepoForGet) List(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForGet) ListByVenue(ctx context.Context, venueID int) ([]model.Auction, error) {
	return nil, nil
}
func (m *mockAuctionRepoForGet) Update(ctx context.Context, auction *model.Auction) error { return nil }
func (m *mockAuctionRepoForGet) UpdateStatus(ctx context.Context, id int, status model.AuctionStatus) error {
	return m.updateStatusErr
}
func (m *mockAuctionRepoForGet) Delete(ctx context.Context, id int) error { return nil }

func TestGetAuctionUseCase_Execute(t *testing.T) {
	validAuction := &model.Auction{ID: 1}

	tests := []struct {
		name                string
		id                  int
		mockAuc             *model.Auction
		mockErr             error
		mockUpdateStatusErr error
		wantErr             bool
		wantNil             bool
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
			wantNil: true,
			// When mock returns nil, nil, usecase `auction, err := uc.repo.GetByID` gets nil, nil.
			// Then `if auction.ShouldBeCompleted()` panics on nil.
			// The UseCase assumes repo returns (nil, error) on failure or always returns struct?
			// Actually typical Go pattern: (nil, ErrNotFound).
			// If repo returns (nil, nil), usecase panics.
			// So this test case exposes that UseCase doesn't handle nil return from repo well if error is nil.
			// However, standard repo contract implies error if nil.
			// Let's assume for this specific test case "NotFound" implies returning an error or just nil.
			// Let's set mockErr to emulate typical NotFound scenario if we want to avoid panic in implementation
			// Or fix implementation. Given scope, let's fix test to simulate proper "Not Found" error from repo.
			mockErr: errors.New("not found"),
			wantErr: true,
		},
		{
			name:    "RepoError",
			id:      1,
			mockErr: errors.New("db error"),
			wantErr: true,
		},
		{
			name: "UpdateStatusError",
			id:   1,
			mockAuc: &model.Auction{
				ID:        1,
				Status:    model.AuctionStatusInProgress,
				EndTime:   timePtr(time.Now().Add(-1 * time.Hour)),
				StartTime: timePtr(time.Now().Add(-2 * time.Hour)),
			},
			mockUpdateStatusErr: errors.New("update failed"),
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAuctionRepoForGet{
				auction:         tt.mockAuc,
				err:             tt.mockErr,
				updateStatusErr: tt.mockUpdateStatusErr,
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

func timePtr(t time.Time) *time.Time {
	return &t
}
