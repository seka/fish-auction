package model

import (
	"time"
)

// AuctionPeriod encapsulates the start and end timestamps of an auction.
type AuctionPeriod struct {
	StartAt *time.Time
	EndAt   *time.Time
}

// NewAuctionPeriod creates a new AuctionPeriod.
func NewAuctionPeriod(start, end *time.Time) AuctionPeriod {
	return AuctionPeriod{
		StartAt: start,
		EndAt:   end,
	}
}

// HasTimeRange reports whether both start and end times are present.
func (p AuctionPeriod) HasTimeRange() bool {
	return p.StartAt != nil && p.EndAt != nil
}

// IsBiddingOpen checks if bidding is currently allowed.
func (p AuctionPeriod) IsBiddingOpen(now time.Time) bool {
	if p.StartAt == nil || p.EndAt == nil {
		return false
	}

	return (now.Equal(*p.StartAt) || now.After(*p.StartAt)) && (now.Equal(*p.EndAt) || now.Before(*p.EndAt))
}

// ShouldExtend checks if the auction should be extended given a new bid.
func (p AuctionPeriod) ShouldExtend(now time.Time, threshold time.Duration) bool {
	if p.EndAt == nil {
		return false
	}

	return p.EndAt.Sub(now) <= threshold
}

// Extend returns a new AuctionPeriod with an extended end time.
func (p AuctionPeriod) Extend(duration time.Duration) AuctionPeriod {
	if p.EndAt == nil {
		return p
	}

	newEnd := p.EndAt.Add(duration)
	newP := p
	newP.EndAt = &newEnd
	return newP
}
