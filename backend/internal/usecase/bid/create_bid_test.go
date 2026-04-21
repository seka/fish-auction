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

func intPtr(i int) *int {
	return &i
}

func TestCreateBidUseCase_Execute(t *testing.T) {
	createBidErr := errors.New("create bid failed")
	txErr := errors.New("tx error")
	dbErr := errors.New("db error")

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	fixedNow := time.Date(2024, 1, 1, 10, 0, 0, 0, jst)
	today := time.Date(2024, 1, 1, 0, 0, 0, 0, jst)
	mockClock := mock.NewMockClock(fixedNow)

	tests := []struct {
		name             string
		input            *model.Bid
		listItemsErr     error
		itemFound        bool
		mockItem         *model.AuctionItem
		getAuctionErr    error
		createErr        error
		txErr            error
		wantID           int
		wantErr          error
		wantCreateCalled bool
		wantTxCalled     bool
		wantUpdateCalled bool
		wantNotification bool
		mockAuction      *model.Auction
		buyerFound       bool
		buyerRepoErr     error
		notificationErr  error
	}{
		{
			name: "Success",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			buyerFound:       true,
			itemFound:        true,
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(today, nil, nil),
				Status:  model.AuctionStatusInProgress,
			},
		},
		{
			name: "Success_BidIncrementOK",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1500),
			},
			buyerFound:       true,
			itemFound:        true,
			mockItem:         &model.AuctionItem{ID: 1, AuctionID: 1, HighestBid: bpp(1000)},
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(today, nil, nil),
				Status:  model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_BidTooLow",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1499),
			},
			buyerFound:       true,
			itemFound:        true,
			mockItem:         &model.AuctionItem{ID: 1, AuctionID: 1, HighestBid: bpp(1000)},
			wantErr:          &domainErrors.ValidationError{Field: "price"},
			mockAuction:      nil,
			wantCreateCalled: false,
			wantTxCalled:     true,
		},
		{
			name: "Error_BuyerNotFound",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			buyerFound:       false,
			wantErr:          &domainErrors.ForbiddenError{},
			wantCreateCalled: false,
			wantTxCalled:     true,
		},
		{
			name: "Error_AuctionStatusNotInProgress",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			buyerFound: true,
			itemFound:  true,
			mockAuction: &model.Auction{
				ID:     1,
				Status: model.AuctionStatusScheduled,
			},
			wantErr:          &domainErrors.ConflictError{},
			wantCreateCalled: false,
			wantTxCalled:     true,
		},
		{
			name: "Success_AutomaticExtension",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			buyerFound:       true,
			itemFound:        true,
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantUpdateCalled: true,
			mockAuction: func() *model.Auction {
				startTime := fixedNow.Add(-1 * time.Hour)
				endTime := fixedNow.Add(2 * time.Minute)

				return &model.Auction{
					ID:      1,
					VenueID: 1,
					Period:  model.NewAuctionPeriod(today, &startTime, &endTime),
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
			buyerFound:       true,
			itemFound:        true,
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantUpdateCalled: false,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(today, nil, nil),
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
			buyerFound:       true,
			itemFound:        true,
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantUpdateCalled: false,
			mockAuction: func() *model.Auction {
				endTime := fixedNow.Add(2 * time.Minute)

				return &model.Auction{
					ID:      1,
					VenueID: 1,
					Period:  model.NewAuctionPeriod(today, nil, &endTime),
					Status:  model.AuctionStatusInProgress,
				}
			}(),
		},
		{
			name: "Error_CreateBidFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			buyerFound:       true,
			itemFound:        true,
			createErr:        createBidErr,
			wantErr:          createBidErr,
			wantCreateCalled: true,
			wantTxCalled:     true,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(today, nil, nil),
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
			buyerFound:   true,
			itemFound:    true,
			txErr:        txErr,
			wantErr:      txErr,
			wantTxCalled: true,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(today, nil, nil),
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
			buyerFound: true,
			itemFound:  true,
			wantErr:    &domainErrors.ValidationError{Field: "auction_time"},
			mockAuction: func() *model.Auction {
				startTime := today.Add(12 * time.Hour)
				endTime := today.Add(13 * time.Hour)
				return &model.Auction{
					ID:     1,
					Period: model.NewAuctionPeriod(today, &startTime, &endTime),
					Status: model.AuctionStatusInProgress,
				}
			}(),
			wantTxCalled: true,
		},
		{
			name:         "Error_ItemRepoError",
			buyerFound:   true,
			input:        &model.Bid{ItemID: 1, BuyerID: 1, Price: bp(1000)},
			listItemsErr: dbErr,
			wantErr:      dbErr,
			wantTxCalled: true,
		},
		{
			name:         "Error_ItemNotFound",
			buyerFound:   true,
			input:        &model.Bid{ItemID: 1, BuyerID: 1, Price: bp(1000)},
			itemFound:    false,
			listItemsErr: nil,
			wantErr:      &domainErrors.NotFoundError{Resource: "Item"},
			wantTxCalled: true,
		},
		{
			name:          "Error_AuctionRepoError",
			buyerFound:    true,
			input:         &model.Bid{ItemID: 1, BuyerID: 1, Price: bp(1000)},
			itemFound:     true,
			getAuctionErr: dbErr,
			wantErr:       dbErr,
			wantTxCalled:  true,
		},
		{
			name:         "Error_AuctionNotFound",
			buyerFound:   true,
			input:        &model.Bid{ItemID: 1, BuyerID: 1, Price: bp(1000)},
			itemFound:    true,
			mockAuction:  nil,
			wantErr:      &domainErrors.NotFoundError{Resource: "Auction"},
			wantTxCalled: true,
		},
		{
			name: "Success_SendNotification_WhenOutbid",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 2,
				Price:   bp(2000),
			},
			buyerFound: true,
			itemFound:  true,
			mockItem: &model.AuctionItem{
				ID:              1,
				AuctionID:       1,
				FishType:        "Maguro",
				HighestBid:      bpp(1500),
				HighestBidderID: intPtr(1),
			},
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantNotification: true,
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(today, nil, nil),
				Status:  model.AuctionStatusInProgress,
			},
		},
		{
			name: "Error_BuyerRepoError",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(1000),
			},
			buyerFound:   true,
			buyerRepoErr: dbErr,
			txErr:        nil,
			wantErr:      dbErr,
			wantTxCalled: true,
		},
		{
			name: "Success_NoNotification_WhenSameBidder",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   bp(2000),
			},
			buyerFound: true,
			itemFound:  true,
			mockItem: &model.AuctionItem{
				ID:              1,
				AuctionID:       1,
				FishType:        "Maguro",
				HighestBid:      bpp(1500),
				HighestBidderID: intPtr(1),
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
		{
			name: "Success_EvenIfNotificationFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 2,
				Price:   bp(2000),
			},
			buyerFound: true,
			itemFound:  true,
			mockItem: &model.AuctionItem{
				ID:              1,
				AuctionID:       1,
				FishType:        "Maguro",
				HighestBid:      bpp(1500),
				HighestBidderID: intPtr(1),
			},
			wantID:           1,
			wantCreateCalled: true,
			wantTxCalled:     true,
			wantNotification: true,
			notificationErr:  errors.New("notification failed"),
			mockAuction: &model.Auction{
				ID:      1,
				VenueID: 1,
				Period:  model.NewAuctionPeriod(today, nil, nil),
				Status:  model.AuctionStatusInProgress,
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

			mockBuyerRepo := &mock.MockBuyerRepository{
				FindByIDFunc: func(_ context.Context, _ int) (*model.Buyer, error) {
					if tt.buyerRepoErr != nil {
						return nil, tt.buyerRepoErr
					}
					if !tt.buyerFound {
						return nil, nil
					}
					return &model.Buyer{ID: tt.input.BuyerID}, nil
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
			mockPushUseCase := &mock.MockPublishNotificationUseCase{
				ExecuteFunc: func(_ context.Context, _ int, _ any) error {
					notificationCalled = true
					return tt.notificationErr
				},
			}

			mockCacheInv := &mock.MockCacheInvalidator{
				InvalidateCacheFunc: func(_ context.Context, _ int) error {
					return nil
				},
			}

			uc := bid.NewCreateBidUseCase(mockItemRepo, mockBuyerRepo, mockBidRepo, mockAuctionRepo, mockPushUseCase, mockTxMgr, mockCacheInv, mockClock)
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
					var wantForbErr *domainErrors.ForbiddenError
					if errors.As(tt.wantErr, &wantForbErr) {
						var gotForbErr *domainErrors.ForbiddenError
						if !errors.As(err, &gotForbErr) {
							t.Fatalf("expected ForbiddenError, got %T: %v", err, err)
						}
					} else {
						var wantNotFoundErr *domainErrors.NotFoundError
						if errors.As(tt.wantErr, &wantNotFoundErr) {
							var gotNotFoundErr *domainErrors.NotFoundError
							if !errors.As(err, &gotNotFoundErr) {
								t.Fatalf("expected NotFoundError, got %T: %v", err, err)
							}
							if wantNotFoundErr.Resource != "" && gotNotFoundErr.Resource != wantNotFoundErr.Resource {
								t.Fatalf("expected resource %s, got %s", wantNotFoundErr.Resource, gotNotFoundErr.Resource)
							}
						} else {
							var wantConflictErr *domainErrors.ConflictError
							if errors.As(tt.wantErr, &wantConflictErr) {
								var gotConflictErr *domainErrors.ConflictError
								if !errors.As(err, &gotConflictErr) {
									t.Fatalf("expected ConflictError, got %T: %v", err, err)
								}
							} else if !errors.Is(err, tt.wantErr) {
								t.Fatalf("expected error %v, got %v", tt.wantErr, err)
							}
						}
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
