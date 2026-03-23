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

func bp(amount int) model.BidPrice {
	return model.NewBidPrice(amount)
}

func bpp(amount int) *model.BidPrice {
	p := model.NewBidPrice(amount)
	return &p
}

//go:fix inline

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
		mockItem         *model.AuctionItem // Custom item for specific tests (e.g. HighestBid)
		getAuctionErr    error
		updateStatusErr  error
		createErr        error
		txErr            error
		wantID           int
		wantErr          error
		wantCreateCalled bool
		wantTxCalled     bool
		wantUpdateCalled bool
		wantNotification bool // Check if SendNotification is called
		mockAuction      *model.Auction
	}{
		{
			name: "Success",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			itemFound:        true,
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(time.Now(), nil, nil),
				Status:  model.AuctionStatusInProgress,
			},
		},
		{
			name: "Success_BidIncrementOK",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1500), // 1000 + 500 (min increment for < 10000 is 500)
			},
			itemFound:        true,
			mockItem:         &model.AuctionItem{ID: 1, AuctionID: 1, HighestBid: bpp(1000)}, // Current price 1000
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(time.Now(), nil, nil),
				Status:  model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_BidTooLow",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1499), // 1000 + 500 = 1500 required
			},
			itemFound:        true,
			mockItem:         &model.AuctionItem{ID: 1, AuctionID: 1, HighestBid: bpp(1000)},
			wantErr:          &domainErrors.ValidationError{Field: "price"},
			mockAuction:      nil,
			wantCreateCalled: false,
			wantTxCalled:     false,
		},
		{
			name: "Success_AutomaticExtension",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			itemFound:        true,
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantUpdateCalled: true,
			mockAuction: func() *model.Auction {
				jst := time.FixedZone("Asia/Tokyo", 9*60*60)
				now := time.Now().In(jst)
				startTime := now.Add(-1 * time.Hour)
				endTime := now.Add(2 * time.Minute) // Within 5 mins -> Trigger extension

				return &model.Auction{
					ID:      1,
					VenueID: 1,
					Period:  model.NewAuctionPeriod(now, &startTime, &endTime),
					Status:  model.AuctionStatusInProgress,
				}
			}(),
		},
		{
			name: "Success_NoExtensionNeeded",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			itemFound:        true,
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantUpdateCalled: false,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(time.Now(), nil, nil),
				Status:  model.AuctionStatusInProgress,
			},
		},
		{
			name: "Success_MissingStartTimeSkipsAuctionWindowChecks",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			itemFound:        true,
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantUpdateCalled: false,
			mockAuction: func() *model.Auction {
				jst := time.FixedZone("Asia/Tokyo", 9*60*60)
				now := time.Now().In(jst)
				endTime := now.Add(2 * time.Minute)

				return &model.Auction{
					ID:      1,
					VenueID: 1,
					Period:  model.NewAuctionPeriod(now, nil, &endTime),
					Status:  model.AuctionStatusInProgress,
				}
			}(),
		},
		{
			name: "Error_UpdateStatusFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			itemFound:        true,
			updateStatusErr:  updateStatusErr,
			wantErr:          nil,
			wantTxCalled:     true,
			wantCreateCalled: true,
			wantID:           1,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(time.Now(), nil, nil),
				Status:  model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_CreateBidFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			itemFound:        true,
			createErr:        createBidErr,
			wantErr:          createBidErr,
			wantCreateCalled: true,
			wantTxCalled:     true,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(time.Now(), nil, nil),
				Status:  model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_TransactionManagerFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			itemFound:    true,
			txErr:        txErr,
			wantErr:      txErr,
			wantTxCalled: true,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(time.Now(), nil, nil),
				Status:  model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_AuctionPeriodInvalid",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			itemFound: true,
			wantErr:   &domainErrors.ValidationError{Field: "auction_time"},
			mockAuction: func() *model.Auction {
				jst := time.FixedZone("Asia/Tokyo", 9*60*60)
				now := time.Now().In(jst)
				startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
				endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 1, 0, 0, jst)
				if now.Hour() == 0 && now.Minute() < 2 {
					startTime = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, jst)
					endTime = time.Date(now.Year(), now.Month(), now.Day(), 12, 1, 0, 0, jst)
				}
				return &model.Auction{
					ID:     1,
					Period: model.NewAuctionPeriod(now, &startTime, &endTime),
					Status: model.AuctionStatusInProgress,
				}
			}(),
		},
		{
			name:         "Error_ItemRepoError",
			input:        &model.Bid{ItemID: 1, BuyerID: 1, Price: bp(1000)},
			listItemsErr: dbErr,
			wantErr:      dbErr,
		},
		{
			name:         "Error_ItemNotFound",
			input:        &model.Bid{ItemID: 1, BuyerID: 1, Price: bp(1000)},
			itemFound:    false,
			listItemsErr: nil,
			wantErr:      &domainErrors.ValidationError{Field: "item_id"},
		},
		{
			name:          "Error_AuctionRepoError",
			input:         &model.Bid{ItemID: 1, BuyerID: 1, Price: bp(1000)},
			itemFound:     true,
			getAuctionErr: dbErr,
			wantErr:       dbErr,
		},
		{
			name: "Success_SendNotification_WhenOutbid",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 2, // New bidder
				Price:   bp(2000),
			},
			itemFound: true,
			mockItem: &model.AuctionItem{
				ID:              1,
				AuctionID:       1,
				FishType:        "Maguro",
				HighestBid:      bpp(1500),
				HighestBidderID: new(1),
			},
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantNotification: true,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(time.Now(), nil, nil),
				Status:  model.AuctionStatusInProgress,
			},
		},
		{
			name: "Success_NoNotification_WhenSameBidder",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1, // Same bidder
				Price:   bp(2000),
			},
			itemFound: true,
			mockItem: &model.AuctionItem{
				ID:              1,
				FishType:        "Maguro",
				HighestBid:      bpp(1500),
				HighestBidderID: new(1),
			},
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantNotification: false,
			mockAuction: &model.Auction{
				ID:     1,
				Status: model.AuctionStatusInProgress,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCalled := false
			createCalled := false
			txCalled := false

			mockItemRepo := &mock.MockItemRepository{
				FindByIDWithLockFunc: func(_ context.Context, _ int) (*model.AuctionItem, error) {
					if tt.listItemsErr != nil {
						return nil, tt.listItemsErr
					}
					if !tt.itemFound {
						return nil, nil
					}
					if tt.mockItem != nil {
						return tt.mockItem, nil
					}
					return &model.AuctionItem{
						ID:        tt.input.ItemID,
						AuctionID: 1,
						FishType:  "Aji",
					}, nil
				},
			}

			mockBidRepo := &mock.MockBidRepository{
				CreateFunc: func(_ context.Context, b *model.Bid) (*model.Bid, error) {
					createCalled = true
					if b.ItemID != tt.input.ItemID || b.BuyerID != tt.input.BuyerID || b.Price.Amount() != tt.input.Price.Amount() {
						t.Fatalf("bid field mismatch: got %+v, want %+v", b, tt.input)
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
				FindByIDWithLockFunc: func(_ context.Context, _ int) (*model.Auction, error) {
					if tt.getAuctionErr != nil {
						return nil, tt.getAuctionErr
					}
					return tt.mockAuction, nil
				},
				UpdateFunc: func(_ context.Context, _ *model.Auction) error {
					updateCalled = true
					return nil
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

			notificationCalled := false
			mockPushUseCase := &mock.MockPushNotificationUseCase{
				SendNotificationFunc: func(_ context.Context, _ int, _ any) error {
					notificationCalled = true
					return nil
				},
			}

			mockCacheInv := &mock.MockCacheInvalidator{
				InvalidateCacheFunc: func(_ context.Context, _ int) error {
					return nil
				},
			}

			uc := bid.NewCreateBidUseCase(mockItemRepo, mockBidRepo, mockAuctionRepo, mockPushUseCase, mockTxMgr, mockCacheInv)
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
				} else if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if created != nil {
					t.Fatalf("expected nil result, got %+v", created)
				}
			} else if tt.wantErr == nil { // No error expected
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if created == nil || created.ID != tt.wantID {
					t.Fatalf("unexpected created bid %+v", created)
				}
			}

			if notificationCalled != tt.wantNotification {
				t.Fatalf("SendNotification called = %v, want %v", notificationCalled, tt.wantNotification)
			}
			if updateCalled != tt.wantUpdateCalled {
				t.Fatalf("Update called = %v, want %v", updateCalled, tt.wantUpdateCalled)
			}
			if createCalled != tt.wantCreateCalled {
				t.Fatalf("Create called = %v, want %v", createCalled, tt.wantCreateCalled)
			}
			if !txCalled {
				t.Fatalf("WithTransaction was not called")
			}
		})
	}
}
