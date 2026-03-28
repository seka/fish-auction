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

// ShouldBeCompleted checks if the auction should be completed based on time
func (a *Auction) ShouldBeCompleted() bool {
	if a.Status == AuctionStatusCompleted || a.Status == AuctionStatusCancelled {
		return false
	}

	tz := NewTimeZone(LocationJST)
	now := tz.Now()
	location := tz.Location()

	// If auction date is in the past
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	auctionDateInTZ := tz.At(a.Period.AuctionDate)
	auctionDate := time.Date(
		auctionDateInTZ.Year(),
		auctionDateInTZ.Month(),
		auctionDateInTZ.Day(),
		0,
		0,
		0,
		0,
		location,
	)
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
