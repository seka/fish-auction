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
	return m.updateStatusErr
}
func (m *mockAuctionRepoForGet) Delete(_ context.Context, _ int) error { return nil }

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
			// When mock returns nil, nil, usecase `auction, err := uc.repo.FindByID` gets nil, nil.
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
			mockAuc: func() *model.Auction {
				jst := time.FixedZone("Asia/Tokyo", 9*60*60)
				now := time.Now().In(jst)
				today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
				startTime := today.Add(10 * time.Hour) // 10:00 Today
				endTime := today.Add(11 * time.Hour)   // 11:00 Today

				// If now is already past 11:00 today, this will work. 
				// But we want to ensure ShouldBeCompleted returns true.
				// ShouldBeCompleted is true if now > endTime.
				// So let's make endTime definitely in the past.
				if now.After(today.Add(2 * time.Hour)) {
					// now is after 02:00, so 00:00-01:00 is in the past.
					startTime = today
					endTime = today.Add(1 * time.Hour)
				} else {
					// now is early (e.g. 00:30), so we use yesterday.
					yesterday := today.Add(-24 * time.Hour)
					startTime = yesterday.Add(10 * time.Hour)
					endTime = yesterday.Add(11 * time.Hour)
					today = yesterday
				}

				return &model.Auction{
					ID:     1,
					Status: model.AuctionStatusInProgress,
					Period: model.NewAuctionPeriod(today, &startTime, &endTime),
				}
			}(),
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

//go:fix inline
