package dto

import "time"

// Auction DTOs
type CreateAuctionRequest struct {
	VenueID     int     `json:"venue_id"`
	AuctionDate string  `json:"auction_date"` // YYYY-MM-DD
	StartTime   *string `json:"start_time"`   // HH:MM:SS
	EndTime     *string `json:"end_time"`     // HH:MM:SS
	Status      string  `json:"status"`
}

type UpdateAuctionRequest struct {
	VenueID     int     `json:"venue_id"`
	AuctionDate string  `json:"auction_date"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Status      string  `json:"status"`
}

type UpdateAuctionStatusRequest struct {
	Status string `json:"status"`
}

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
