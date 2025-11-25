package usecase

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// BidUseCase defines the interface for bidding-related business logic
type BidUseCase interface {
	Bid(ctx context.Context, t *model.Transaction) (*model.Transaction, error)
}

type bidInteractor struct {
	itemRepo repository.ItemRepository
	txRepo   repository.TransactionRepository
}

func NewBidInteractor(itemRepo repository.ItemRepository, txRepo repository.TransactionRepository) BidUseCase {
	return &bidInteractor{itemRepo: itemRepo, txRepo: txRepo}
}

func (i *bidInteractor) Bid(ctx context.Context, t *model.Transaction) (*model.Transaction, error) {
	// Note: Ideally, this should be transactional.
	// For now, we keep the logic simple as per previous implementation, but separated.
	// In a real Clean Architecture, we might need a UnitOfWork pattern or similar for transactions across repositories.

	if err := i.itemRepo.UpdateStatus(ctx, t.ItemID, "Sold"); err != nil {
		return nil, err
	}

	return i.txRepo.Create(ctx, t)
}
