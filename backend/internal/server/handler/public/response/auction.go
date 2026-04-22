package response

import "time"

// Auction represents a public view of an auction.
type Auction struct {
	ID        int       `json:"id"`
	VenueID   int       `json:"venue_id"`
	StartAt   *string   `json:"start_at"`
	EndAt     *string   `json:"end_at"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
