package dto

import "time"


// CreateAuctionRequest is a data transfer object.
type CreateAuctionRequest struct {
	VenueID     int     `json:"venue_id"`
	AuctionDate string  `json:"auction_date"` // YYYY-MM-DD
	StartTime   *string `json:"start_time"`   // HH:MM:SS
	EndTime     *string `json:"end_time"`     // HH:MM:SS
	Status      string  `json:"status"`
}

// UpdateAuctionRequest represents the request body for updating an auction.
type UpdateAuctionRequest struct {
	VenueID     int     `json:"venue_id"`
	AuctionDate string  `json:"auction_date"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Status      string  `json:"status"`
}

// UpdateAuctionStatusRequest represents the request body for updating an auction's status.
type UpdateAuctionStatusRequest struct {
	Status string `json:"status"`
}

// AuctionResponse represents the response body for an auction.
type AuctionResponse struct {
	ID          int       `json:"id"`
	VenueID     int       `json:"venue_id"`
	AuctionDate string    `json:"auction_date"`
	StartTime   *string   `json:"start_time"`
	EndTime     *string   `json:"end_time"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
