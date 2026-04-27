package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuction_ShouldBeCompleted(t *testing.T) {
	jst := NewTimeZone(LocationJST).Location()

	// Fixed reference time: 2024-01-01 12:00:00 JST
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, jst)

	newAuction := func(status AuctionStatus, startAt, endAt *time.Time) *Auction {
		return &Auction{
			Status: status,
			Period: NewAuctionPeriod(startAt, endAt),
		}
	}

	yesterday9am := time.Date(2023, 12, 31, 9, 0, 0, 0, jst)
	tomorrow9am := time.Date(2024, 1, 2, 9, 0, 0, 0, jst)
	today9am := time.Date(2024, 1, 1, 9, 0, 0, 0, jst)
	today11am := time.Date(2024, 1, 1, 11, 0, 0, 0, jst)
	today13pm := time.Date(2024, 1, 1, 13, 0, 0, 0, jst)

	tests := []struct {
		name    string
		status  AuctionStatus
		startAt *time.Time
		endAt   *time.Time
		want    bool
	}{
		{
			name:    "returns false when status is already completed",
			status:  AuctionStatusCompleted,
			startAt: &yesterday9am,
			want:    false,
		},
		{
			name:    "returns false when status is canceled",
			status:  AuctionStatusCancelled,
			startAt: &yesterday9am,
			want:    false,
		},
		{
			name:    "returns true when auction date is in the past",
			status:  AuctionStatusInProgress,
			startAt: &yesterday9am,
			want:    true,
		},
		{
			name:    "returns false when auction date is in the future",
			status:  AuctionStatusInProgress,
			startAt: &tomorrow9am,
			want:    false,
		},
		{
			name:    "returns true when today's auction end time has passed",
			status:  AuctionStatusInProgress,
			startAt: &today9am,
			endAt:   &today11am,
			want:    true,
		},
		{
			name:    "returns false when today's auction end time is still ahead",
			status:  AuctionStatusInProgress,
			startAt: &today9am,
			endAt:   &today13pm,
			want:    false,
		},
		{
			name:    "returns false when today's auction has no end time",
			status:  AuctionStatusInProgress,
			startAt: &today9am,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auction := newAuction(tt.status, tt.startAt, tt.endAt)
			assert.Equal(t, tt.want, auction.ShouldBeCompleted(now))
		})
	}
}

//go:fix inline
