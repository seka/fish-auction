package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateBidUseCase defines the interface for creating bids
type CreateBidUseCase interface {
	Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error)
}

// createBidUseCase handles the creation of bids
type createBidUseCase struct {
	itemRepo    repository.ItemRepository
	bidRepo     repository.BidRepository
	auctionRepo repository.AuctionRepository
	txMgr       repository.TransactionManager
}

// NewCreateBidUseCase creates a new instance of CreateBidUseCase
func NewCreateBidUseCase(
	itemRepo repository.ItemRepository,
	bidRepo repository.BidRepository,
	auctionRepo repository.AuctionRepository,
	txMgr repository.TransactionManager,
) CreateBidUseCase {
	return &createBidUseCase{
		itemRepo:    itemRepo,
		bidRepo:     bidRepo,
		auctionRepo: auctionRepo,
		txMgr:       txMgr,
	}
}

// Execute creates a new bid
func (uc *createBidUseCase) Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	// Get item to find auction_id
	items, err := uc.itemRepo.List(ctx, "")
	if err != nil {
		return nil, err
	}

	var item *model.AuctionItem
	for i := range items {
		if items[i].ID == bid.ItemID {
			item = &items[i]
			break
		}
	}

	if item == nil {
		return nil, &errors.ValidationError{
			Field:   "item_id",
			Message: "item not found",
		}
	}

	// Get auction to check bidding period
	auction, err := uc.auctionRepo.GetByID(ctx, item.AuctionID)
	if err != nil {
		return nil, err
	}

	// Check if auction is within bidding hours
	if auction.StartTime != nil && auction.EndTime != nil {
		now := time.Now()

		// Create start and end datetime
		startDateTime := time.Date(
			auction.AuctionDate.Year(), auction.AuctionDate.Month(), auction.AuctionDate.Day(),
			auction.StartTime.Hour(), auction.StartTime.Minute(), auction.StartTime.Second(), 0, now.Location(),
		)
		endDateTime := time.Date(
			auction.AuctionDate.Year(), auction.AuctionDate.Month(), auction.AuctionDate.Day(),
			auction.EndTime.Hour(), auction.EndTime.Minute(), auction.EndTime.Second(), 0, now.Location(),
		)

		if now.Before(startDateTime) || now.After(endDateTime) {
			return nil, &errors.ValidationError{
				Field: "auction_time",
				Message: fmt.Sprintf("Bidding is not allowed outside auction hours (%02d:%02d - %02d:%02d)",
					auction.StartTime.Hour(), auction.StartTime.Minute(),
					auction.EndTime.Hour(), auction.EndTime.Minute()),
			}
		}
	}

	var result *model.Bid

	err = uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// Create bid record
		created, err := uc.bidRepo.Create(txCtx, bid)
		if err != nil {
			return err
		}

		result = created
		return nil
	})

	return result, err
}
