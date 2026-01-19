package dto

import "time"

// Item DTOs
type CreateItemRequest struct {
	AuctionID   int    `json:"auction_id"`
	FishermanID int    `json:"fisherman_id"`
	FishType    string `json:"fish_type"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
}

type UpdateItemRequest struct {
	AuctionID   int    `json:"auction_id"`
	FishermanID int    `json:"fisherman_id"`
	FishType    string `json:"fish_type"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
	Status      string `json:"status"`
}

type UpdateItemSortOrderRequest struct {
	SortOrder int `json:"sort_order"`
}

type ItemResponse struct {
	ID                int       `json:"id"`
	AuctionID         int       `json:"auction_id"`
	FishermanID       int       `json:"fisherman_id"`
	FishType          string    `json:"fish_type"`
	Quantity          int       `json:"quantity"`
	Unit              string    `json:"unit"`
	Status            string    `json:"status"`
	HighestBid        *int      `json:"highest_bid,omitempty"`
	HighestBidderID   *int      `json:"highest_bidder_id,omitempty"`
	HighestBidderName *string   `json:"highest_bidder_name,omitempty"`
	SortOrder         int       `json:"sort_order"`
	CreatedAt         time.Time `json:"created_at"`
}
