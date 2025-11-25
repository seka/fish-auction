package entity

import (
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type Fisherman struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (e *Fisherman) ToModel() *model.Fisherman {
	return &model.Fisherman{
		ID:   e.ID,
		Name: e.Name,
	}
}

type Buyer struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (e *Buyer) ToModel() *model.Buyer {
	return &model.Buyer{
		ID:   e.ID,
		Name: e.Name,
	}
}

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
