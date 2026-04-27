package response

// Auction represents an auction view for the buyer.
type Auction struct {
	ID        int     `json:"id"`
	VenueID   int     `json:"venue_id"`
	StartAt   *string `json:"start_at"`
	EndAt     *string `json:"end_at"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
