package response

import "time"

// Auction represents a public view of an auction.
type Auction struct {
	ID          int       `json:"id"`
	VenueID     int       `json:"venue_id"`
	AuctionDate string    `json:"auction_date"`
	StartTime   *string   `json:"start_time"`
	EndTime     *string   `json:"end_time"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
