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
		auction, err := u.auctionRepo.FindByID(txCtx, item.AuctionID)
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

		// 5. Validate bid amount
		currentAmount := 0
		if item.HighestBid != nil {
			currentAmount = item.HighestBid.Amount()
		}
		if bid.Price.Amount() <= currentAmount {
			return &domainErrors.ValidationError{Field: "price", Message: "Bid amount must be greater than current price"}
		}

		// 6. Create bid
		bid.CreatedAt = now
		createdBid, err = u.bidRepo.Create(txCtx, bid)
		if err != nil {
			return fmt.Errorf("failed to create bid: %w", err)
		}

		// 7. Update item highest bid
		previousHighestBidderID := item.HighestBidderID
		previousAmount := currentAmount

		item.HighestBid = &bid.Price
		item.HighestBidderID = &bid.BuyerID
		if _, err := u.itemRepo.Update(txCtx, item); err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}

		// 8. Automatic Extension
		if auction.Period.ShouldExtend(now, 5*time.Minute) {
			auction.Period = auction.Period.Extend(5 * time.Minute)
			if err := u.auctionRepo.Update(txCtx, auction); err != nil {
				return fmt.Errorf("failed to extend auction: %w", err)
			}
		}

		// 9. Notify outbid buyer
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
	payload := map[string]interface{}{
		"type":            "outbid",
		"item_id":         item.ID,
		"fish_type":       item.FishType,
		"previous_amount": previousAmount,
		"current_amount":  item.HighestBid.Amount(),
	}
	return u.outboxRepo.InsertPushNotificationJob(ctx, buyerID, payload)
}
