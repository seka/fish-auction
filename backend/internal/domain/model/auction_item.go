package model

import "time"

type AuctionItem struct {
	ID          int
	FishermanID int
	FishType    string
	Quantity    int
	Unit        string
	Status      ItemStatus
	CreatedAt   time.Time
}
