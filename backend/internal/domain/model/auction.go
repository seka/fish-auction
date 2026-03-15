package model

import (
	"time"
)

// AuctionStatus represents the status of an auction
type AuctionStatus string

const (
	AuctionStatusScheduled  AuctionStatus = "scheduled"
	AuctionStatusInProgress AuctionStatus = "in_progress"
	AuctionStatusCompleted  AuctionStatus = "completed"
	AuctionStatusCancelled  AuctionStatus = "canceled"
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

// ShouldBeCompleted checks if the auction should be completed based on time
func (a *Auction) ShouldBeCompleted() bool {
	if a.Status == AuctionStatusCompleted || a.Status == AuctionStatusCancelled {
		return false
	}

	tz := NewTimeZone(LocationJST)
	now := tz.Now()

	// If auction date is in the past
	today := now.Truncate(24 * time.Hour)
	auctionDate := a.Period.AuctionDate.Truncate(24 * time.Hour)
	if auctionDate.Before(today) {
		return true
	}

	// If auction date is today, check end time using AuctionPeriod method
	if auctionDate.Equal(today) {
		end := a.Period.GetEndDateTime()
		if end != nil && now.After(*end) {
			return true
		}
	}

	return false
}
