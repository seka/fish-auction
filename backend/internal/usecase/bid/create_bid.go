package bid

import (
	"context"
	"fmt"
	"time"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

const (
	// AuctionExtensionThreshold は終了直前のこの時間内に入札が入った場合に
	// 自動延長を行う閾値。
	AuctionExtensionThreshold = 5 * time.Minute

	// AuctionExtensionDuration は自動延長で延ばす時間。
	AuctionExtensionDuration = 5 * time.Minute
)

// CreateBidUseCase defines the interface for creating a bid.
type CreateBidUseCase interface {
	// Execute creates a new bid and updates the item's current price.
	Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error)
}

type createBidUseCase struct {
	itemRepo     repository.ItemRepository
	buyerRepo    repository.BuyerRepository
	bidRepo      repository.BidRepository
	auctionRepo  repository.AuctionRepository
	outboxRepo   repository.OutboxRepository
	txMgr        repository.TransactionManager
	itemCacheInv repository.CacheInvalidator
	clock        service.Clock
}

var _ CreateBidUseCase = (*createBidUseCase)(nil)

// NewCreateBidUseCase creates a new instance of CreateBidUseCase.
func NewCreateBidUseCase(
	itemRepo repository.ItemRepository,
	buyerRepo repository.BuyerRepository,
	bidRepo repository.BidRepository,
	auctionRepo repository.AuctionRepository,
	outboxRepo repository.OutboxRepository,
	txMgr repository.TransactionManager,
	itemCacheInv repository.CacheInvalidator,
	clock service.Clock,
) CreateBidUseCase {
	return &createBidUseCase{
		itemRepo:     itemRepo,
		buyerRepo:    buyerRepo,
		bidRepo:      bidRepo,
		auctionRepo:  auctionRepo,
		outboxRepo:   outboxRepo,
		txMgr:        txMgr,
		itemCacheInv: itemCacheInv,
		clock:        clock,
	}
}

func (u *createBidUseCase) Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	var createdBid *model.Bid
	err := u.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Verify buyer exists
		buyer, err := u.buyerRepo.FindByID(txCtx, bid.BuyerID)
		if err != nil {
			return fmt.Errorf("failed to verify buyer: %w", err)
		}
		if buyer == nil {
			return &domainErrors.ForbiddenError{Message: "Buyer not found"}
		}

		// 2. Get and lock item
		item, err := u.itemRepo.FindByIDWithLock(txCtx, bid.ItemID)
		if err != nil {
			return fmt.Errorf("failed to find item: %w", err)
		}
		if item == nil {
			return &domainErrors.NotFoundError{Resource: "Item", ID: bid.ItemID}
		}

		// 3. Get auction and validate status
		// 自動延長で auction.Period を更新する可能性があるため、行ロックを取得して
		// 同一商品への並行入札による Period の lost update を防ぐ。
		auction, err := u.auctionRepo.FindByIDWithLock(txCtx, item.AuctionID)
		if err != nil {
			return fmt.Errorf("failed to find auction: %w", err)
		}
		if auction == nil {
			return &domainErrors.NotFoundError{Resource: "Auction", ID: item.AuctionID}
		}
		if auction.Status != model.AuctionStatusInProgress {
			return &domainErrors.ConflictError{Message: "Auction is not in progress"}
		}

		// 4. Validate bid time
		now := u.clock.Now()
		if !auction.Period.IsBiddingOpen(now) {
			return &domainErrors.ValidationError{Field: "auction_time", Message: "Bid is outside of auction period"}
		}

		// 5. Validate bid amount with minimum increment
		currentPrice := model.NewBidPrice(0)
		if item.HighestBid != nil {
			currentPrice = *item.HighestBid
		}
		minIncrement := currentPrice.CalculateMinIncrement()
		minAcceptable := currentPrice.Add(minIncrement)
		if bid.Price.LessThan(minAcceptable) {
			return &domainErrors.ValidationError{
				Field:   "price",
				Message: fmt.Sprintf("Bid price must be at least %d", minAcceptable.Amount()),
			}
		}

		// 6. Create bid
		// item.HighestBid / HighestBidderID は transactions テーブルから都度算出される
		// derived 値であり、auction_items テーブルには永続化しないため明示的な Update は不要。
		previousHighestBidderID := item.HighestBidderID
		previousAmount := currentPrice.Amount()

		bid.CreatedAt = now
		createdBid, err = u.bidRepo.Create(txCtx, bid)
		if err != nil {
			return fmt.Errorf("failed to create bid: %w", err)
		}

		// 7. Automatic Extension
		if auction.Period.ShouldExtend(now, AuctionExtensionThreshold) {
			auction.Period = auction.Period.Extend(AuctionExtensionDuration)
			if err := u.auctionRepo.Update(txCtx, auction); err != nil {
				return fmt.Errorf("failed to extend auction: %w", err)
			}
		}

		// 8. Notify outbid buyer
		if previousHighestBidderID != nil && *previousHighestBidderID != bid.BuyerID {
			if err := u.notifyOutbid(txCtx, item, *previousHighestBidderID, previousAmount); err != nil {
				fmt.Printf("failed to enqueue outbid notification: %v\n", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 10. Invalidate cache
	if err := u.itemCacheInv.InvalidateCache(ctx, bid.ItemID); err != nil {
		fmt.Printf("failed to invalidate item cache: %v\n", err)
	}

	return createdBid, nil
}

func (u *createBidUseCase) notifyOutbid(ctx context.Context, item *model.AuctionItem, buyerID, previousAmount int) error {
	title := "高値更新"
	body := fmt.Sprintf("%s への入札が更新されました（¥%d → ¥%d）", item.FishType, previousAmount, item.HighestBid.Amount())
	// フロントは個別商品ページを持たず、商品はオークション詳細ページ (/auctions/[id]) で一覧表示される。
	url := fmt.Sprintf("/auctions/%d", item.AuctionID)
	return u.outboxRepo.InsertPushJob(ctx, model.JobTypePushOutbid, buyerID, title, body, url)
}
