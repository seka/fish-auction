package model

import "time"

// AuctionStatus represents the status of an auction
type AuctionStatus string

const (
	AuctionStatusScheduled  AuctionStatus = "scheduled"
	AuctionStatusInProgress AuctionStatus = "in_progress"
	AuctionStatusCompleted  AuctionStatus = "completed"
	AuctionStatusCancelled  AuctionStatus = "cancelled"
)

// IsValid checks if the auction status is valid
func (s AuctionStatus) IsValid() bool {
	switch s {
	case AuctionStatusScheduled, AuctionStatusInProgress, AuctionStatusCompleted, AuctionStatusCancelled:
		return true
	default:
		return false
	}
}

// Auction represents an auction event (セリイベント)
type Auction struct {
	ID          int
	VenueID     int
	AuctionDate time.Time
	StartTime   *time.Time
	EndTime     *time.Time
	Status      AuctionStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
