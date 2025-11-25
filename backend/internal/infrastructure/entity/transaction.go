package entity

import (
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type Transaction struct {
	ID        int       `db:"id"`
	ItemID    int       `db:"item_id"`
	BuyerID   int       `db:"buyer_id"`
	Price     int       `db:"price"`
	CreatedAt time.Time `db:"created_at"`
}

func (e *Transaction) ToModel() *model.Transaction {
	return &model.Transaction{
		ID:        e.ID,
		ItemID:    e.ItemID,
		BuyerID:   e.BuyerID,
		Price:     e.Price,
		CreatedAt: e.CreatedAt,
	}
}
