package response

// Auction represents an auction view for the buyer.
type Auction struct {
	ID          int     `json:"id"`
	VenueID     int     `json:"venue_id"`
	AuctionDate string  `json:"auction_date"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
