package model

import "time"

type AuctionItem struct {
	ID          int       `json:"id"`
	FishermanID int       `json:"fisherman_id"`
	FishType    string    `json:"fish_type"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
