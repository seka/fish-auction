package auction

import (
	"context"
	"log"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/usecase/notification"
)

// UpdateAuctionStatusUseCase defines the interface for updating auction status.
type UpdateAuctionStatusUseCase interface {
	// Execute updates an auction's status and potentially sends notifications.
	Execute(ctx context.Context, id int, status model.AuctionStatus) error
}

type updateAuctionStatusUseCase struct {
	auctionRepo                repository.AuctionRepository
	buyerRepo                  repository.BuyerRepository
	publishNotificationUseCase notification.PublishNotificationUseCase
}

var _ UpdateAuctionStatusUseCase = (*updateAuctionStatusUseCase)(nil)

func NewUpdateAuctionStatusUseCase(
	auctionRepo repository.AuctionRepository,
	buyerRepo repository.BuyerRepository,
	publishNotification notification.PublishNotificationUseCase,
) UpdateAuctionStatusUseCase {
	return &updateAuctionStatusUseCase{
		auctionRepo:                auctionRepo,
		buyerRepo:                  buyerRepo,
		publishNotificationUseCase: publishNotification,
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

			payload := map[string]any{
				"title": title,
				"body":  body,
				"url":   "/auctions",
			}

			for _, b := range buyers {
				if err := uc.publishNotificationUseCase.Execute(ctx, b.ID, payload); err != nil {
					// 通知失敗はログ出力のみ行い、全体の処理に影響を与えない
					log.Printf("failed to send notification for auction status update: %v", err)
				}
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
