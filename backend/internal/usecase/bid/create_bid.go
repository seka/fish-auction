package bid

import (
	"context"
	"fmt"
	"log"

	"time"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/usecase/notification"
)

const (
	// AuctionExtensionThreshold is the time remaining before the auction ends
	// during which a new bid will trigger an extension.
	AuctionExtensionThreshold = 5 * time.Minute

	// AuctionExtensionDuration is the duration by which the auction will be
	// extended when ShouldExtend is true.
	AuctionExtensionDuration = 5 * time.Minute
)

// CreateBidUseCase defines the interface for creating a bid.
type CreateBidUseCase interface {
	// Execute creates a new bid.
	Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error)
}

type createBidUseCase struct {
	itemRepo                   repository.ItemRepository
	buyerRepo                  repository.BuyerRepository
	bidRepo                    repository.BidRepository
	auctionRepo                repository.AuctionRepository
	publishNotificationUseCase notification.PublishNotificationUseCase
	txMgr                      repository.TransactionManager
	itemCacheInv               repository.CacheInvalidator
	clock                      service.Clock
}

var _ CreateBidUseCase = (*createBidUseCase)(nil)

// NewCreateBidUseCase creates a new instance of CreateBidUseCase.
func NewCreateBidUseCase(
	itemRepo repository.ItemRepository,
	buyerRepo repository.BuyerRepository,
	bidRepo repository.BidRepository,
	auctionRepo repository.AuctionRepository,
	publishNotificationUseCase notification.PublishNotificationUseCase,
	txMgr repository.TransactionManager,
	itemCacheInv repository.CacheInvalidator,
	clock service.Clock,
) CreateBidUseCase {
	return &createBidUseCase{
		itemRepo:                   itemRepo,
		buyerRepo:                  buyerRepo,
		bidRepo:                    bidRepo,
		auctionRepo:                auctionRepo,
		publishNotificationUseCase: publishNotificationUseCase,
		txMgr:                      txMgr,
		itemCacheInv:               itemCacheInv,
		clock:                      clock,
	}
}

// Execute creates a new bid.
func (uc *createBidUseCase) Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	var result *model.Bid
	var lockedItem *model.AuctionItem

	err := uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. 購入者の存在確認
		if err := uc.verifyBuyer(txCtx, bid.BuyerID); err != nil {
			return err
		}

		// 2. 最新の商品情報を取得（ロック付）して検証
		item, err := uc.getAndValidateItem(txCtx, bid.ItemID, bid.Price)
		if err != nil {
			return err
		}
		lockedItem = item

		// 3. オークション情報を取得（ロック付）して検証
		auction, err := uc.getAndValidateAuction(txCtx, item.AuctionID)
		if err != nil {
			return err
		}

		// 4. 自動延長処理
		if err := uc.extendAuctionIfNeeded(txCtx, auction); err != nil {
			return err
		}

		// 5. 入札レコードを作成
		created, err := uc.bidRepo.Create(txCtx, bid)
		if err != nil {
			return err
		}

		result = created
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 入札成功後の後処理（トランザクション外）
	uc.invalidateCache(ctx, result.ItemID)
	uc.notifyOutbid(ctx, result, lockedItem)

	return result, nil
}

func (uc *createBidUseCase) verifyBuyer(ctx context.Context, buyerID int) error {
	buyer, err := uc.buyerRepo.FindByID(ctx, buyerID)
	if err != nil {
		return err
	}
	if buyer == nil {
		return &errors.ForbiddenError{
			Message: "bidder record not found for authenticated user",
		}
	}
	return nil
}

func (uc *createBidUseCase) getAndValidateItem(ctx context.Context, itemID int, price model.BidPrice) (*model.AuctionItem, error) {
	item, err := uc.itemRepo.FindByIDWithLock(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, &errors.NotFoundError{
			Resource: "Item",
			ID:       itemID,
		}
	}

	if err := uc.validateBidPrice(item, price); err != nil {
		return nil, err
	}

	return item, nil
}

func (uc *createBidUseCase) getAndValidateAuction(ctx context.Context, auctionID int) (*model.Auction, error) {
	auction, err := uc.auctionRepo.FindByIDWithLock(ctx, auctionID)
	if err != nil {
		return nil, err
	}
	if auction == nil {
		return nil, &errors.NotFoundError{
			Resource: "Auction",
			ID:       auctionID,
		}
	}

	if auction.Status != model.AuctionStatusInProgress {
		return nil, &errors.ConflictError{
			Message: fmt.Sprintf("bidding is not allowed for auction with status %s", auction.Status),
		}
	}

	if err := uc.validateAuctionPeriod(auction); err != nil {
		return nil, err
	}

	return auction, nil
}

func (uc *createBidUseCase) validateBidPrice(item *model.AuctionItem, bidPrice model.BidPrice) error {
	currentPrice := model.NewBidPrice(0)
	if item.HighestBid != nil {
		currentPrice = *item.HighestBid
	}

	minIncrement := currentPrice.CalculateMinIncrement()
	if bidPrice.LessThan(currentPrice.Add(minIncrement)) {
		return &errors.ValidationError{
			Field:   "price",
			Message: fmt.Sprintf("Bid price must be at least %d", currentPrice.Add(minIncrement).Amount()),
		}
	}
	return nil
}

func (uc *createBidUseCase) validateAuctionPeriod(auction *model.Auction) error {
	if !auction.Period.HasTimeRange() {
		return nil
	}

	now := uc.clock.NowIn(model.LocationJST)
	if !auction.Period.IsBiddingOpen(now) {
		start := auction.Period.GetStartDateTime()
		end := auction.Period.GetEndDateTime()
		return &errors.ValidationError{
			Field: "auction_time",
			Message: fmt.Sprintf("Bidding is not allowed outside auction hours (%02d:%02d - %02d:%02d)",
				start.Hour(), start.Minute(),
				end.Hour(), end.Minute()),
		}
	}
	return nil
}

func (uc *createBidUseCase) extendAuctionIfNeeded(ctx context.Context, auction *model.Auction) error {
	if !auction.Period.HasTimeRange() {
		return nil
	}

	now := uc.clock.NowIn(model.LocationJST)

	if auction.Period.ShouldExtend(now, AuctionExtensionThreshold) {
		auction.Period = auction.Period.Extend(AuctionExtensionDuration)

		if err := uc.auctionRepo.Update(ctx, auction); err != nil {
			return fmt.Errorf("failed to extend auction: %w", err)
		}
	}

	return nil
}

func (uc *createBidUseCase) invalidateCache(ctx context.Context, itemID int) {
	_ = uc.itemCacheInv.InvalidateCache(ctx, itemID)
}

func (uc *createBidUseCase) notifyOutbid(ctx context.Context, bid *model.Bid, item *model.AuctionItem) {
	if item.HighestBidderID != nil && *item.HighestBidderID != bid.BuyerID {
		log.Printf("Outbid detected. Sending notification to previous bidder (ID: %d). Current bidder (ID: %d)", *item.HighestBidderID, bid.BuyerID)
		payload := map[string]any{
			"title": "高値更新されました",
			"body":  fmt.Sprintf("%s の価格が %d 円に更新されました。", item.FishType, bid.Price.Amount()),
			"url":   fmt.Sprintf("/auctions/%d", item.AuctionID),
		}
		if err := uc.publishNotificationUseCase.Execute(ctx, *item.HighestBidderID, payload); err != nil {
			// 通知失敗はログ出力のみ行い、全体の処理に影響を与えない
			log.Printf("failed to send notification for outbid: %v", err)
		}
	}
}
