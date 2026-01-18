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

	// Helper to create int pointer
	intPtr := func(i int) *int { return &i }

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
			name: "Success_BidIncrementOK",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1500, // 1000 + 500 (min increment for < 10000 is 500, wait. 1000 < 10000 -> 500)
			},
			itemFound:        true,
			mockItem:         &model.AuctionItem{ID: 1, AuctionID: 1, HighestBid: intPtr(1000)}, // Current price 1000
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			mockAuction: &model.Auction{
				ID:          1,
				VenueID:     1,
				AuctionDate: time.Now(),
				StartTime:   nil,
				EndTime:     nil, // No time check if nil
				Status:      model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_BidTooLow",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1499, // 1000 + 500 = 1500 required
			},
			itemFound:   true,
			mockItem:    &model.AuctionItem{ID: 1, AuctionID: 1, HighestBid: intPtr(1000)},
			wantErr:     &domainErrors.ValidationError{Field: "price"},
			mockAuction: nil, // Should fail before fetching auction? No, after item check. But item check is first?
			// Logic: Get Item -> Validate Price -> Get Auction. So Auction might not be fetched if price invalid.
			// Re-reading code: Yes.
			wantCreateCalled: false,
			wantTxCalled:     false,
		},
		{
			name: "Success_AutomaticExtension",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
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
				// 日付跨ぎ対策: もしendTimeが翌日になってしまったら、AuctionDateも調整するか、
				// そもそもAuctionDateをnowの日付に合わせる（これは既にされているが、time.Dateで時分秒だけ入れるロジックが問題）

				return &model.Auction{
					ID:          1,
					VenueID:     1,
					AuctionDate: now, // nowの日付を使用
					StartTime:   &startTime,
					EndTime:     &endTime,
					Status:      model.AuctionStatusInProgress,
				}
			}(),
		},
		{
			name: "Success_NoExtensionNeeded",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			itemFound:        true,
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantUpdateCalled: false,
			mockAuction: &model.Auction{
				ID:          1,
				VenueID:     1,
				AuctionDate: time.Now(),
				StartTime:   nil, // 時刻チェックをスキップ
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
			wantErr:          nil,
			wantTxCalled:     true,
			wantCreateCalled: true,
			wantID:           1,
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
				jst := time.FixedZone("Asia/Tokyo", 9*60*60)
				now := time.Now().In(jst)
				// 明らかに過去の時間を設定して、現在時刻が「時間外」になるようにする
				// ただし、実施日に依存しないよう AuctionDate は今日にする
				startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
				endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 1, 0, 0, jst)
				// もし今が 00:00 〜 00:01 の間なら 12:00 等にずらす
				if now.Hour() == 0 && now.Minute() < 2 {
					startTime = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, jst)
					endTime = time.Date(now.Year(), now.Month(), now.Day(), 12, 1, 0, 0, jst)
				}
				return &model.Auction{
					ID:          1,
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
		{
			name: "Success_SendNotification_WhenOutbid",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 2, // New bidder
				Price:   2000,
			},
			itemFound: true,
			mockItem: &model.AuctionItem{
				ID:              1,
				AuctionID:       1,
				FishType:        "Maguro",
				HighestBid:      intPtr(1500),
				HighestBidderID: intPtr(1), // Previous bidder (different from current)
			},
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantNotification: true,
			mockAuction: &model.Auction{
				ID:          1,
				VenueID:     1,
				AuctionDate: time.Now(),
				Status:      model.AuctionStatusInProgress,
			},
		},
		{
			name: "Success_NoNotification_WhenSameBidder",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1, // Same bidder
				Price:   2000,
			},
			itemFound: true,
			mockItem: &model.AuctionItem{
				ID:              1,
				FishType:        "Maguro",
				HighestBid:      intPtr(1500),
				HighestBidderID: intPtr(1), // Same bidder
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
				FindByIDFunc: func(ctx context.Context, id int) (*model.AuctionItem, error) {
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
				InvalidateCacheFunc: func(ctx context.Context, id int) error {
					return nil
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
					// Only return mockAuction if it's set.
					// If getAuctionErr is set, this might not be reached depending on implementation logic order.
					// Tests above with getAuctionErr set expect error.
					return tt.mockAuction, nil
				},
				UpdateFunc: func(ctx context.Context, auction *model.Auction) error {
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
				SendNotificationFunc: func(ctx context.Context, buyerID int, payload interface{}) error {
					notificationCalled = true
					return nil
				},
			}

			uc := bid.NewCreateBidUseCase(mockItemRepo, mockBidRepo, mockAuctionRepo, mockPushUseCase, mockTxMgr)
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

			if notificationCalled != tt.wantNotification {
				t.Fatalf("SendNotification called = %v, want %v", notificationCalled, tt.wantNotification)
			}
			if updateCalled != tt.wantUpdateCalled {
				t.Fatalf("Update called = %v, want %v", updateCalled, tt.wantUpdateCalled)
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
