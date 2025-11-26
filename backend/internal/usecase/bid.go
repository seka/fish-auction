package usecase

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// BidUseCase defines the interface for bidding-related business logic
type BidUseCase interface {
	Bid(ctx context.Context, bid *model.Bid) (*model.Bid, error)
}

type bidInteractor struct {
	itemRepo repository.ItemRepository
	bidRepo  repository.BidRepository
	txMgr    repository.TransactionManager
}

func NewBidInteractor(itemRepo repository.ItemRepository, bidRepo repository.BidRepository, txMgr repository.TransactionManager) BidUseCase {
	return &bidInteractor{
		itemRepo: itemRepo,
		bidRepo:  bidRepo,
		txMgr:    txMgr,
	}
}

func (i *bidInteractor) Bid(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	var result *model.Bid

	err := i.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// Update item status to "Sold"
		if err := i.itemRepo.UpdateStatus(txCtx, bid.ItemID, model.ItemStatusSold); err != nil {
			return err
		}

		// Create bid record
		created, err := i.bidRepo.Create(txCtx, bid)
		if err != nil {
			return err
		}

		result = created
		return nil
	})

	return result, err
}
