package request

// CreateAuction holds data for auction creation.
type CreateAuction struct {
	VenueID     int     `json:"venue_id"`
	AuctionDate string  `json:"auction_date"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Status      string  `json:"status"`
}

// UpdateAuction holds data for updating an auction.
type UpdateAuction struct {
	VenueID     int     `json:"venue_id"`
	AuctionDate string  `json:"auction_date"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Status      string  `json:"status"`
}

// UpdateAuctionStatus holds data for updating an auction's status.
type UpdateAuctionStatus struct {
	Status string `json:"status"`
}
