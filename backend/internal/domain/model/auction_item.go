package model

import "time"

type AuctionItem struct {
	ID          int
	AuctionID   int
	FishermanID int
	FishType    string
	Quantity    int
	Unit        string
	Status      ItemStatus
	HighestBid  *int
	CreatedAt   time.Time
}
