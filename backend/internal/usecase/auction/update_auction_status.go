package auction

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase/notification"
)

// UpdateAuctionStatusUseCase defines the interface for updating auction status
type UpdateAuctionStatusUseCase interface {
	Execute(ctx context.Context, id int, status model.AuctionStatus) error
}

// updateAuctionStatusUseCase handles updating auction status
type updateAuctionStatusUseCase struct {
	auctionRepo repository.AuctionRepository
	buyerRepo   repository.BuyerRepository
	pushUseCase notification.PushNotificationUseCase
}

// NewUpdateAuctionStatusUseCase creates a new instance of UpdateAuctionStatusUseCase
func NewUpdateAuctionStatusUseCase(
	auctionRepo repository.AuctionRepository,
	buyerRepo repository.BuyerRepository,
	pushUseCase notification.PushNotificationUseCase,
) UpdateAuctionStatusUseCase {
	return &updateAuctionStatusUseCase{
		auctionRepo: auctionRepo,
		buyerRepo:   buyerRepo,
		pushUseCase: pushUseCase,
	}
}

// Execute updates an auction's status
func (uc *updateAuctionStatusUseCase) Execute(ctx context.Context, id int, status model.AuctionStatus) error {
	// Validate status
	if !status.IsValid() {
		return &InvalidStatusError{Status: string(status)}
	}

	if err := uc.auctionRepo.UpdateStatus(ctx, id, status); err != nil {
		return err
	}

	// オークションが開始または終了した際に全買付人に通知
	if status == model.AuctionStatusInProgress || status == model.AuctionStatusCompleted {
		buyers, err := uc.buyerRepo.List(ctx)
		if err == nil {
			title := "オークション開始"
			body := "新しいオークションが開始されました。"
			if status == model.AuctionStatusCompleted {
				title = "オークション終了"
				body = "オークションが終了しました。結果をご確認ください。"
			}

			payload := map[string]interface{}{
				"title": title,
				"body":  body,
				"url":   "/auctions",
			}

			for _, b := range buyers {
				_ = uc.pushUseCase.SendNotification(ctx, b.ID, payload)
			}
		}
	}

	return nil
}

// InvalidStatusError represents an invalid auction status error
type InvalidStatusError struct {
	Status string
}

func (e *InvalidStatusError) Error() string {
	return "invalid auction status: " + e.Status
}
