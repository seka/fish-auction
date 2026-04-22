package model

import (
	"time"
)

// AuctionStatus represents the status of an auction
type AuctionStatus string

const (
	// AuctionStatusScheduled provides AuctionStatusScheduled related functionality.
	AuctionStatusScheduled AuctionStatus = "scheduled"
	// AuctionStatusInProgress provides AuctionStatusInProgress related functionality.
	AuctionStatusInProgress AuctionStatus = "in_progress"
	// AuctionStatusCompleted provides AuctionStatusCompleted related functionality.
	AuctionStatusCompleted AuctionStatus = "completed"
	// AuctionStatusCancelled represents a auctionstatuscancelled in the system.
	AuctionStatusCancelled AuctionStatus = "canceled"
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
	ID        int
	VenueID   int
	Period    AuctionPeriod
	Status    AuctionStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ShouldBeCompleted checks if the auction should be completed based on the provided time
func (a *Auction) ShouldBeCompleted(now time.Time) bool {
	if a.Status == AuctionStatusCompleted || a.Status == AuctionStatusCancelled {
		return false
	}

	if a.Period.StartAt == nil {
		return false
	}

	location := NewTimeZone(LocationJST).Location()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	startInJST := a.Period.StartAt.In(location)
	startDate := time.Date(startInJST.Year(), startInJST.Month(), startInJST.Day(), 0, 0, 0, 0, location)

	if startDate.Before(today) {
		return true
	}

	if startDate.Equal(today) && a.Period.EndAt != nil && now.After(*a.Period.EndAt) {
		return true
	}

	return false
}
