package repository

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type ItemRepository interface {
	Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
	List(ctx context.Context, status string) ([]model.AuctionItem, error)
	ListByAuction(ctx context.Context, auctionID int) ([]model.AuctionItem, error)
	FindByID(ctx context.Context, id int) (*model.AuctionItem, error)
	Update(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error)
	Delete(ctx context.Context, id int) error
	UpdateStatus(ctx context.Context, id int, status model.ItemStatus) error
	UpdateSortOrder(ctx context.Context, id int, sortOrder int) error
	Reorder(ctx context.Context, auctionID int, ids []int) error
	InvalidateCache(ctx context.Context, id int) error
}
