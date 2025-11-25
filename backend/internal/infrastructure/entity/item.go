package entity

import (
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type AuctionItem struct {
	ID          int       `db:"id"`
	FishermanID int       `db:"fisherman_id"`
	FishType    string    `db:"fish_type"`
	Quantity    int       `db:"quantity"`
	Unit        string    `db:"unit"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}

func (e *AuctionItem) ToModel() *model.AuctionItem {
	return &model.AuctionItem{
		ID:          e.ID,
		FishermanID: e.FishermanID,
		FishType:    e.FishType,
		Quantity:    e.Quantity,
		Unit:        e.Unit,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
	}
}
