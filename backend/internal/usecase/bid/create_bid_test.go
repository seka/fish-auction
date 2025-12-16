package bid_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/bid"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestCreateBidUseCase_Execute(t *testing.T) {
	updateStatusErr := errors.New("update status failed")
	createBidErr := errors.New("create bid failed")
	txErr := errors.New("tx error")
	dbErr := errors.New("db error")

	tests := []struct {
		name             string
		input            *model.Bid
		listItemsErr     error
		itemFound        bool
		getAuctionErr    error
		updateStatusErr  error
		createErr        error
		txErr            error
		wantID           int
		wantErr          error
		wantCreateCalled bool
		wantTxCalled     bool
		mockAuction      *model.Auction
	}{
		{
			name: "Success",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			itemFound:        true,
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			mockAuction: &model.Auction{
				ID:          1,
				VenueID:     1,
				AuctionDate: time.Now(),
				StartTime:   nil,
				EndTime:     nil,
				Status:      model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_UpdateStatusFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			itemFound:        true,
			updateStatusErr:  updateStatusErr,
			wantErr:          nil, // UpdateStatus is no longer called, so no error expected from it
			wantTxCalled:     true,
			wantCreateCalled: true, // Should proceed to create
			wantID:           1,    // Should succeed
			mockAuction: &model.Auction{
				ID:          1,
				VenueID:     1,
				AuctionDate: time.Now(),
				StartTime:   nil,
				EndTime:     nil,
				Status:      model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_CreateBidFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			itemFound:        true,
			createErr:        createBidErr,
			wantErr:          createBidErr,
			wantCreateCalled: true,
			wantTxCalled:     true,
			mockAuction: &model.Auction{
				ID:          1,
				VenueID:     1,
				AuctionDate: time.Now(),
				StartTime:   nil,
				EndTime:     nil,
				Status:      model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_TransactionManagerFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			itemFound:    true,
			txErr:        txErr,
			wantErr:      txErr,
			wantTxCalled: true,
			mockAuction: &model.Auction{
				ID:          1,
				VenueID:     1,
				AuctionDate: time.Now(),
				StartTime:   nil,
				EndTime:     nil,
				Status:      model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_AuctionPeriodInvalid",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			itemFound: true,
			wantErr:   &domainErrors.ValidationError{Field: "auction_time"},
			mockAuction: func() *model.Auction {
				now := time.Now()
				startTime := now.Add(-2 * time.Hour)
				endTime := now.Add(-1 * time.Hour)
				return &model.Auction{
					ID:          1,
					VenueID:     1,
					AuctionDate: now,
					StartTime:   &startTime,
					EndTime:     &endTime,
					Status:      model.AuctionStatusInProgress,
				}
			}(),
		},
		{
			name: "Error_AuctionPeriodInvalid_BeforeStart",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			itemFound: true,
			wantErr:   &domainErrors.ValidationError{Field: "auction_time"},
			mockAuction: func() *model.Auction {
				now := time.Now()
				startTime := now.Add(1 * time.Hour)
				endTime := now.Add(2 * time.Hour)
				return &model.Auction{
					ID:          1,
					VenueID:     1,
					AuctionDate: now,
					StartTime:   &startTime,
					EndTime:     &endTime,
					Status:      model.AuctionStatusInProgress,
				}
			}(),
		},
		{
			name:         "Error_ItemRepoError",
			input:        &model.Bid{ItemID: 1, BuyerID: 1, Price: 1000},
			listItemsErr: dbErr,
			wantErr:      dbErr,
		},
		{
			name:         "Error_ItemNotFound",
			input:        &model.Bid{ItemID: 1, BuyerID: 1, Price: 1000},
			itemFound:    false, // Simulate item not found
			listItemsErr: nil,
			wantErr:      &domainErrors.ValidationError{Field: "item_id"},
		},
		{
			name:          "Error_AuctionRepoError",
			input:         &model.Bid{ItemID: 1, BuyerID: 1, Price: 1000},
			itemFound:     true,
			getAuctionErr: dbErr,
			wantErr:       dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCalled := false
			createCalled := false
			txCalled := false

			mockItemRepo := &mock.MockItemRepository{
				ListFunc: func(ctx context.Context, status string) ([]model.AuctionItem, error) {
					if tt.listItemsErr != nil {
						return nil, tt.listItemsErr
					}
					if !tt.itemFound {
						return []model.AuctionItem{}, nil
					}
					return []model.AuctionItem{
						{ID: tt.input.ItemID, AuctionID: 1},
					}, nil
				},
			}

			mockBidRepo := &mock.MockBidRepository{
				CreateFunc: func(ctx context.Context, b *model.Bid) (*model.Bid, error) {
					createCalled = true
					if b != tt.input {
						t.Fatalf("bid pointer mismatch")
					}
					if tt.createErr != nil {
						return nil, tt.createErr
					}
					cloned := *b
					cloned.ID = tt.wantID
					return &cloned, nil
				},
			}

			mockAuctionRepo := &mock.MockAuctionRepository{
				GetByIDFunc: func(ctx context.Context, id int) (*model.Auction, error) {
					if tt.getAuctionErr != nil {
						return nil, tt.getAuctionErr
					}
					return tt.mockAuction, nil
				},
			}

			mockTxMgr := &mock.MockTransactionManager{
				WithTransactionFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
					txCalled = true
					if tt.txErr != nil {
						return tt.txErr
					}
					return fn(ctx)
				},
			}

			uc := bid.NewCreateBidUseCase(mockItemRepo, mockBidRepo, mockAuctionRepo, mockTxMgr)
			created, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr != nil {
				var wantValErr *domainErrors.ValidationError
				if errors.As(tt.wantErr, &wantValErr) {
					var gotValErr *domainErrors.ValidationError
					if !errors.As(err, &gotValErr) {
						t.Fatalf("expected ValidationError, got %T: %v", err, err)
					}
					if wantValErr.Field != "" && gotValErr.Field != wantValErr.Field {
						t.Fatalf("expected field %s, got %s", wantValErr.Field, gotValErr.Field)
					}
				} else {
					if !errors.Is(err, tt.wantErr) {
						t.Fatalf("expected error %v, got %v", tt.wantErr, err)
					}
				}
				if created != nil {
					t.Fatalf("expected nil result, got %+v", created)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if created == nil || created.ID != tt.wantID {
					t.Fatalf("unexpected created bid %+v", created)
				}
			}

			if updateCalled {
				t.Fatalf("UpdateStatus called = %v, want %v", updateCalled, false)
			}
			if createCalled != tt.wantCreateCalled {
				t.Fatalf("Create called = %v, want %v", createCalled, tt.wantCreateCalled)
			}
			if txCalled != tt.wantTxCalled {
				t.Fatalf("WithTransaction called = %v, want %v", txCalled, tt.wantTxCalled)
			}
		})
	}
}
