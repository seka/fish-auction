package dto

import "time"


// CreateItemRequest is a data transfer object.
type CreateItemRequest struct {
	AuctionID   int    `json:"auction_id"`
	FishermanID int    `json:"fisherman_id"`
	FishType    string `json:"fish_type"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
}

// UpdateItemRequest represents the request body for updating an item.
type UpdateItemRequest struct {
	AuctionID   int    `json:"auction_id"`
	FishermanID int    `json:"fisherman_id"`
	FishType    string `json:"fish_type"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
	Status      string `json:"status"`
}

// UpdateItemSortOrderRequest represents the request body for updating an item's sort order.
type UpdateItemSortOrderRequest struct {
	SortOrder int `json:"sort_order"`
}

// ReorderItemsRequest represents the request body for reordering items.
type ReorderItemsRequest struct {
	IDs []int `json:"ids"`
}

// ItemResponse represents the response body for an item.
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
