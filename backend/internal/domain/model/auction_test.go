package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuction_ShouldBeCompleted(t *testing.T) {
	tz := NewTimeZone(LocationJST)
	location := tz.Location()
	now := tz.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	past := now.Add(-1 * time.Second)
	future := now.Add(1 * time.Second)

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
		skip        bool
		skipReason  string
	}{
		{
			name:        "returns false when status is already completed",
			status:      AuctionStatusCompleted,
			auctionDate: today.AddDate(0, 0, -1),
			want:        false,
		},
		{
			name:        "returns false when status is cancelled",
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
			end:         timePtr(time.Date(0, 1, 1, past.Hour(), past.Minute(), past.Second(), 0, location)),
			want:        true,
			skip:        past.Day() != now.Day() || past.Month() != now.Month() || past.Year() != now.Year(),
			skipReason:  "cannot create an earlier same-day time near midnight",
		},
		{
			name:        "returns false when today's auction end time is still ahead",
			status:      AuctionStatusInProgress,
			auctionDate: today,
			end:         timePtr(time.Date(0, 1, 1, future.Hour(), future.Minute(), future.Second(), 0, location)),
			want:        false,
			skip:        future.Day() != now.Day() || future.Month() != now.Month() || future.Year() != now.Year(),
			skipReason:  "cannot create a later same-day time near midnight",
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
			if tt.skip {
				t.Skip(tt.skipReason)
			}

			auction := newAuction(tt.status, tt.auctionDate, tt.end)
			assert.Equal(t, tt.want, auction.ShouldBeCompleted())
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
