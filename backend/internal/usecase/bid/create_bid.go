package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateBidUseCase は入札作成のインターフェースを定義します
type CreateBidUseCase interface {
	Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error)
}

// createBidUseCase は入札作成処理を扱います
type createBidUseCase struct {
	itemRepo    repository.ItemRepository
	bidRepo     repository.BidRepository
	auctionRepo repository.AuctionRepository
	txMgr       repository.TransactionManager
}

// NewCreateBidUseCase は CreateBidUseCase の新しいインスタンスを作成します
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

// Execute は新しい入札を作成します
func (uc *createBidUseCase) Execute(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	// 商品を取得して auction_id を見つける
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

	// 入札上乗せ額の検証
	currentPrice := 0
	if item.HighestBid != nil {
		currentPrice = *item.HighestBid
	}
	minIncrement := getMinimumBidIncrement(currentPrice)
	if bid.Price < currentPrice+minIncrement {
		return nil, &errors.ValidationError{
			Field:   "price",
			Message: fmt.Sprintf("Bid price must be at least %d", currentPrice+minIncrement),
		}
	}

	// オークション情報を取得して入札期間をチェック
	auction, err := uc.auctionRepo.GetByID(ctx, item.AuctionID)
	if err != nil {
		return nil, err
	}

	// オークションが入札時間内かどうかチェック
	if auction.StartTime != nil && auction.EndTime != nil {
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		now := time.Now().In(jst)

		// 開始日時と終了日時を作成
		startDateTime := time.Date(
			auction.AuctionDate.Year(), auction.AuctionDate.Month(), auction.AuctionDate.Day(),
			auction.StartTime.Hour(), auction.StartTime.Minute(), auction.StartTime.Second(), 0, jst,
		)
		endDateTime := time.Date(
			auction.AuctionDate.Year(), auction.AuctionDate.Month(), auction.AuctionDate.Day(),
			auction.EndTime.Hour(), auction.EndTime.Minute(), auction.EndTime.Second(), 0, jst,
		)

		if now.Before(startDateTime) || now.After(endDateTime) {
			return nil, &errors.ValidationError{
				Field: "auction_time",
				Message: fmt.Sprintf("Bidding is not allowed outside auction hours (%02d:%02d - %02d:%02d)",
					auction.StartTime.Hour(), auction.StartTime.Minute(),
					auction.EndTime.Hour(), auction.EndTime.Minute()),
			}
		}

		// 自動延長: 終了5分前に入札があった場合、5分延長する
		const extensionThreshold = 5 * time.Minute
		const extensionDuration = 5 * time.Minute

		if endDateTime.Sub(now) <= extensionThreshold {
			newEndTime := endDateTime.Add(extensionDuration)
			// 現在の実装では、DBは時間（HH:mm:ss）のみを考慮する設計上の制約がある可能性があります。
			// 本来であれば日付を跨ぐ場合に日付フィールドも更新すべきですが、
			// 現状のMVP/プロトタイプ実装として、終了時間の時間部分のみを更新します。
			// ※ 23:59:59 を超えて延長する場合の挙動については別途検討が必要です。

			// 厳密には、基となる time オブジェクトを更新すべきです。
			newEndTimePure := time.Date(0, 1, 1, newEndTime.Hour(), newEndTime.Minute(), newEndTime.Second(), 0, jst)
			auction.EndTime = &newEndTimePure

			if err := uc.auctionRepo.Update(ctx, auction); err != nil {
				// 延長通知に失敗した場合でも、入札の整合性を保つためエラーとします。
				return nil, fmt.Errorf("failed to extend auction: %w", err)
			}
		}
	}

	var result *model.Bid

	err = uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// 入札レコードを作成
		created, err := uc.bidRepo.Create(txCtx, bid)
		if err != nil {
			return err
		}

		result = created
		return nil
	})

	return result, err
}

func getMinimumBidIncrement(currentPrice int) int {
	if currentPrice < 1000 {
		return 100
	}
	if currentPrice < 10000 {
		return 500
	}
	if currentPrice < 100000 {
		return 1000
	}
	return 5000
}
