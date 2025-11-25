package entity

import (
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type Bid struct {
	ID        int       `db:"id"`
	ItemID    int       `db:"item_id"`
	BuyerID   int       `db:"buyer_id"`
	Price     int       `db:"price"`
	CreatedAt time.Time `db:"created_at"`
}

func (e *Bid) ToModel() *model.Bid {
	return &model.Bid{
		ID:        e.ID,
		ItemID:    e.ItemID,
		BuyerID:   e.BuyerID,
		Price:     e.Price,
		CreatedAt: e.CreatedAt,
	}
}
