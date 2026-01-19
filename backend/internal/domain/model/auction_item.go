package model

import "time"

type AuctionItem struct {
	ID                int
	AuctionID         int
	FishermanID       int
	FishType          string
	Quantity          int
	Unit              string
	Status            ItemStatus
	HighestBid        *int
	HighestBidderID   *int
	HighestBidderName *string
	SortOrder         int
	CreatedAt         time.Time
	DeletedAt         *time.Time
}
