package bid

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateBidUseCase handles the creation of bids with transaction management
type CreateBidUseCase struct {
	itemRepo repository.ItemRepository
	bidRepo  repository.BidRepository
	txMgr    repository.TransactionManager
}

// NewCreateBidUseCase creates a new instance of CreateBidUseCase
func NewCreateBidUseCase(
	itemRepo repository.ItemRepository,
	bidRepo repository.BidRepository,
	txMgr repository.TransactionManager,
) *CreateBidUseCase {
	return &CreateBidUseCase{
		itemRepo: itemRepo,
		bidRepo:  bidRepo,
		txMgr:    txMgr,
	}
}

// Execute creates a new bid and updates the item status atomically
func (uc *CreateBidUseCase) Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	var result *model.Bid

	err := uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// Update item status to "Sold"
		if err := uc.itemRepo.UpdateStatus(txCtx, bid.ItemID, model.ItemStatusSold); err != nil {
			return err
		}

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
