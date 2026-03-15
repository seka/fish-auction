package model

import (
	"time"
)

// AuctionPeriod encapsulates the period (date, start time, and end time) of an auction.
type AuctionPeriod struct {
	AuctionDate time.Time
	StartAt     *time.Time
	EndAt       *time.Time
}

// NewAuctionPeriod creates a new AuctionPeriod.
func NewAuctionPeriod(date time.Time, start, end *time.Time) AuctionPeriod {
	return AuctionPeriod{
		AuctionDate: date,
		StartAt:     start,
		EndAt:       end,
	}
}

// GetStartDateTime returns the absolute start time of the auction.
func (p AuctionPeriod) GetStartDateTime() *time.Time {
	if p.StartAt == nil {
		return nil
	}
	jst := NewTimeZone(LocationJST).Location()
	t := time.Date(
		p.AuctionDate.Year(), p.AuctionDate.Month(), p.AuctionDate.Day(),
		p.StartAt.Hour(), p.StartAt.Minute(), p.StartAt.Second(), 0, jst,
	)
	return &t
}

// GetEndDateTime returns the absolute end time of the auction.
func (p AuctionPeriod) GetEndDateTime() *time.Time {
	if p.EndAt == nil {
		return nil
	}
	jst := NewTimeZone(LocationJST).Location()
	t := time.Date(
		p.AuctionDate.Year(), p.AuctionDate.Month(), p.AuctionDate.Day(),
		p.EndAt.Hour(), p.EndAt.Minute(), p.EndAt.Second(), 0, jst,
	)
	return &t
}

// IsBiddingOpen checks if bidding is currently allowed.
func (p AuctionPeriod) IsBiddingOpen(now time.Time) bool {
	start := p.GetStartDateTime()
	end := p.GetEndDateTime()

	if start == nil || end == nil {
		return false
	}

	return (now.Equal(*start) || now.After(*start)) && (now.Equal(*end) || now.Before(*end))
}

// ShouldExtend checks if the auction should be extended given a new bid.
func (p AuctionPeriod) ShouldExtend(now time.Time, threshold time.Duration) bool {
	end := p.GetEndDateTime()
	if end == nil {
		return false
	}

	return end.Sub(now) <= threshold
}

// Extend returns a new AuctionPeriod with an extended end time.
func (p AuctionPeriod) Extend(duration time.Duration) AuctionPeriod {
	end := p.GetEndDateTime()
	if end == nil {
		return p
	}

	newEnd := end.Add(duration)
	// We only store the time part in EndAt, keeping the original date context.
	// But since EndAt is a *time.Time, we can just use the absolute time,
	// though usually these are stored as column types with only time.
	// For simplicity in the domain model, we just return the new absolute time.
	// When saving to DB, the infra layer will extract the time part.
	jst := NewTimeZone(LocationJST).Location()
	newEndPure := time.Date(0, 1, 1, newEnd.Hour(), newEnd.Minute(), newEnd.Second(), 0, jst)

	newP := p
	newP.EndAt = &newEndPure
	return newP
}
