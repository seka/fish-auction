package bid

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase/notification"
)

// CreateBidUseCase defines the interface for creating a bid.
type CreateBidUseCase interface {
	// Execute creates a new bid.
	Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error)
}

type createBidUseCase struct {
	itemRepo     repository.ItemRepository
	bidRepo      repository.BidRepository
	auctionRepo  repository.AuctionRepository
	pushUseCase  notification.PushNotificationUseCase
	txMgr        repository.TransactionManager
	itemCacheInv repository.CacheInvalidator
}

var _ CreateBidUseCase = (*createBidUseCase)(nil)

// NewCreateBidUseCase creates a new instance of CreateBidUseCase.
func NewCreateBidUseCase(
	itemRepo repository.ItemRepository,
	bidRepo repository.BidRepository,
	auctionRepo repository.AuctionRepository,
	pushUseCase notification.PushNotificationUseCase,
	txMgr repository.TransactionManager,
	itemCacheInv repository.CacheInvalidator,
) *createBidUseCase {
	return &createBidUseCase{
		itemRepo:     itemRepo,
		bidRepo:      bidRepo,
		auctionRepo:  auctionRepo,
		pushUseCase:  pushUseCase,
		txMgr:        txMgr,
		itemCacheInv: itemCacheInv,
	}
}

// Execute creates a new bid.
func (uc *createBidUseCase) Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	var result *model.Bid
	var lockedItem *model.AuctionItem

	err := uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. ロックを取得して最新の商品情報を取得
		item, err := uc.itemRepo.FindByIDWithLock(txCtx, bid.ItemID)
		if err != nil {
			return err
		}
		if item == nil {
			return &errors.ValidationError{
				Field:   "item_id",
				Message: "item not found",
			}
		}
		lockedItem = item

		// 2. 入札価格の検証
		if err := uc.validateBidPrice(item, bid.Price); err != nil {
			return err
		}

		// 3. オークション情報を取得（必要に応じてロック）して入札期間をチェック
		auction, err := uc.auctionRepo.FindByIDWithLock(txCtx, item.AuctionID)
		if err != nil {
			return err
		}

		// 4. 入札期間のチェック
		if err := uc.validateAuctionPeriod(auction); err != nil {
			return err
		}

		// 5. 自動延長処理
		if err := uc.extendAuctionIfNeeded(txCtx, auction); err != nil {
			return err
		}

		// 6. 入札レコードを作成
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

	// 7. 入札成功後の後処理（トランザクション外）
	uc.invalidateCache(ctx, result.ItemID)
	uc.notifyOutbid(ctx, result, lockedItem)

	return result, nil
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
	tz := model.NewTimeZone(model.LocationJST)
	now := tz.Now()
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
	tz := model.NewTimeZone(model.LocationJST)
	now := tz.Now()

	const extensionThreshold = 5 * time.Minute
	const extensionDuration = 5 * time.Minute

	if auction.Period.ShouldExtend(now, extensionThreshold) {
		auction.Period = auction.Period.Extend(extensionDuration)

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
		payload := map[string]interface{}{
			"title": "高値更新されました",
			"body":  fmt.Sprintf("%s の価格が %d 円に更新されました。", item.FishType, bid.Price.Amount()),
			"url":   fmt.Sprintf("/auctions/%d", item.AuctionID),
		}
		_ = uc.pushUseCase.SendNotification(ctx, *item.HighestBidderID, payload)
	}
}
