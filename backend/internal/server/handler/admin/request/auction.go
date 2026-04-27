package request

// CreateAuction holds data for auction creation.
type CreateAuction struct {
	VenueID int     `json:"venue_id"`
	StartAt *string `json:"start_at"`
	EndAt   *string `json:"end_at"`
	Status  string  `json:"status"`
}

// UpdateAuction holds data for updating an auction.
type UpdateAuction struct {
	VenueID int     `json:"venue_id"`
	StartAt *string `json:"start_at"`
	EndAt   *string `json:"end_at"`
	Status  string  `json:"status"`
}

// UpdateAuctionStatus holds data for updating an auction's status.
type UpdateAuctionStatus struct {
	Status  string  `json:"status"`
	StartAt *string `json:"start_at"`
}
