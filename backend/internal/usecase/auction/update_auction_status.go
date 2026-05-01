package auction

import (
	"context"
	"fmt"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// UpdateAuctionStatusUseCase defines the interface for updating an auction's status
type UpdateAuctionStatusUseCase interface {
	// Execute updates an auction's status
	Execute(ctx context.Context, id int, status model.AuctionStatus) error
}

type updateAuctionStatusUseCase struct {
	auctionRepo repository.AuctionRepository
	buyerRepo   repository.BuyerRepository
	outboxRepo  repository.OutboxRepository
	txMgr       repository.TransactionManager
}

var _ UpdateAuctionStatusUseCase = (*updateAuctionStatusUseCase)(nil)

func NewUpdateAuctionStatusUseCase(
	auctionRepo repository.AuctionRepository,
	buyerRepo repository.BuyerRepository,
	outboxRepo repository.OutboxRepository,
	txMgr repository.TransactionManager,
) UpdateAuctionStatusUseCase {
	return &updateAuctionStatusUseCase{
		auctionRepo: auctionRepo,
		buyerRepo:   buyerRepo,
		outboxRepo:  outboxRepo,
		txMgr:       txMgr,
	}
}

// Execute updates an auction's status
func (uc *updateAuctionStatusUseCase) Execute(ctx context.Context, id int, status model.AuctionStatus) error {
	// Validate status
	if !status.IsValid() {
		return &InvalidStatusError{Status: string(status)}
	}

	return uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// Update status
		if err := uc.auctionRepo.UpdateStatus(txCtx, id, status); err != nil {
			return fmt.Errorf("failed to update auction status: %w", err)
		}

		// Notify buyers (in a real app, this might be filtered by subscription)
		buyers, err := uc.buyerRepo.List(txCtx)
		if err != nil {
			return fmt.Errorf("failed to list buyers for notification: %w", err)
		}

		notificationPayload := map[string]interface{}{
			"type":       "auction_status_changed",
			"auction_id": id,
			"new_status": string(status),
		}

		for _, buyer := range buyers {
			if err := uc.outboxRepo.InsertPushNotificationJob(txCtx, buyer.ID, notificationPayload); err != nil {
				// Log error but continue for other buyers
				fmt.Printf("failed to enqueue notification for buyer %d: %v\n", buyer.ID, err)
			}
		}

		return nil
	})
}

// InvalidStatusError is returned when the auction status is invalid.
type InvalidStatusError struct {
	Status string
}

func (e *InvalidStatusError) Error() string {
	return fmt.Sprintf("invalid auction status: %s", e.Status)
}
