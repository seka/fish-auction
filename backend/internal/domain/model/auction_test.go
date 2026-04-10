package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuction_ShouldBeCompleted(t *testing.T) {
	tz := NewTimeZone(LocationJST)
	location := tz.Location()

	// Fixed reference time: 2024-01-01 12:00:00 JST
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, location)
	today := time.Date(2024, 1, 1, 0, 0, 0, 0, location)

	newAuction := func(status AuctionStatus, auctionDate time.Time, end *time.Time) *Auction {
		return &Auction{
			Status: status,
			Period: NewAuctionPeriod(auctionDate, nil, end),
		}
	}

	tests := []struct {
		name        string
		status      AuctionStatus
		auctionDate time.Time
		end         *time.Time
		want        bool
	}{
		{
			name:        "returns false when status is already completed",
			status:      AuctionStatusCompleted,
			auctionDate: today.AddDate(0, 0, -1),
			want:        false,
		},
		{
			name:        "returns false when status is canceled",
			status:      AuctionStatusCancelled,
			auctionDate: today.AddDate(0, 0, -1),
			want:        false,
		},
		{
			name:        "returns true when auction date is in the past",
			status:      AuctionStatusInProgress,
			auctionDate: today.AddDate(0, 0, -1),
			want:        true,
		},
		{
			name:        "returns false when auction date is in the future",
			status:      AuctionStatusInProgress,
			auctionDate: today.AddDate(0, 0, 1),
			want:        false,
		},
		{
			name:        "returns true when today's auction end time has passed",
			status:      AuctionStatusInProgress,
			auctionDate: today,
			end:         func() *time.Time { t := today.Add(11 * time.Hour); return &t }(),
			want:        true,
		},
		{
			name:        "returns false when today's auction end time is still ahead",
			status:      AuctionStatusInProgress,
			auctionDate: today,
			end:         func() *time.Time { t := today.Add(13 * time.Hour); return &t }(),
			want:        false,
		},
		{
			name:        "returns false when today's auction has no end time",
			status:      AuctionStatusInProgress,
			auctionDate: today,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auction := newAuction(tt.status, tt.auctionDate, tt.end)
			assert.Equal(t, tt.want, auction.ShouldBeCompleted(now))
		})
	}
}

//go:fix inline
