package model

import "time"

type Fisherman struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Buyer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AuctionItem struct {
	ID          int       `json:"id"`
	FishermanID int       `json:"fisherman_id"`
	FishType    string    `json:"fish_type"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	Status      string    `json:"status"` // Pending, Sold
	CreatedAt   time.Time `json:"created_at"`
}

type Transaction struct {
	ID        int       `json:"id"`
	ItemID    int       `json:"item_id"`
	BuyerID   int       `json:"buyer_id"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type InvoiceItem struct {
	BuyerID     int    `json:"buyer_id"`
	BuyerName   string `json:"buyer_name"`
	TotalAmount int    `json:"total_amount"` // Price + Tax
}
