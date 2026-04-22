package response

import "time"

// Item represents a public view of an auction item.
type Item struct {
	ID          int       `json:"id"`
	AuctionID   int       `json:"auction_id"`
	FishermanID int       `json:"fisherman_id"`
	FishType    string    `json:"fish_type"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	HighestBid  *int      `json:"highest_bid,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}
