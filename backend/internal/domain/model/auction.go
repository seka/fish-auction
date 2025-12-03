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

// ShouldBeCompleted checks if the auction should be completed based on time
func (a *Auction) ShouldBeCompleted() bool {
	if a.Status == AuctionStatusCompleted || a.Status == AuctionStatusCancelled {
		return false
	}

	// Use JST
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)

	// Truncate auction date to day (00:00:00)
	// Treat AuctionDate as JST date
	auctionDate := time.Date(a.AuctionDate.Year(), a.AuctionDate.Month(), a.AuctionDate.Day(), 0, 0, 0, 0, jst)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)

	// If auction date is in the past
	if auctionDate.Before(today) {
		return true
	}

	// If auction date is today, check end time
	if auctionDate.Equal(today) && a.EndTime != nil {
		endDateTime := time.Date(
			auctionDate.Year(), auctionDate.Month(), auctionDate.Day(),
			a.EndTime.Hour(), a.EndTime.Minute(), a.EndTime.Second(), 0, jst,
		)
		if now.After(endDateTime) {
			return true
		}
	}

	return false
}
