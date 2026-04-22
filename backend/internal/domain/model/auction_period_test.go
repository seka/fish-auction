package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuctionPeriod_IsBiddingOpen(t *testing.T) {
	jst := NewTimeZone(LocationJST).Location()
	start := time.Date(2026, 3, 15, 9, 0, 0, 0, jst)
	end := time.Date(2026, 3, 15, 17, 0, 0, 0, jst)

	p := NewAuctionPeriod(&start, &end)

	tests := []struct {
		name     string
		now      time.Time
		expected bool
	}{
		{
			name:     "before start",
			now:      time.Date(2026, 3, 15, 8, 59, 59, 0, jst),
			expected: false,
		},
		{
			name:     "exactly at start",
			now:      time.Date(2026, 3, 15, 9, 0, 0, 0, jst),
			expected: true,
		},
		{
			name:     "during period",
			now:      time.Date(2026, 3, 15, 12, 0, 0, 0, jst),
			expected: true,
		},
		{
			name:     "exactly at end",
			now:      time.Date(2026, 3, 15, 17, 0, 0, 0, jst),
			expected: true,
		},
		{
			name:     "after end",
			now:      time.Date(2026, 3, 15, 17, 0, 1, 0, jst),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, p.IsBiddingOpen(tt.now))
		})
	}
}

func TestAuctionPeriod_ShouldExtend(t *testing.T) {
	jst := NewTimeZone(LocationJST).Location()
	end := time.Date(2026, 3, 15, 17, 0, 0, 0, jst)
	p := NewAuctionPeriod(nil, &end)

	threshold := 5 * time.Minute

	tests := []struct {
		name     string
		now      time.Time
		expected bool
	}{
		{
			name:     "well before end",
			now:      time.Date(2026, 3, 15, 16, 54, 0, 0, jst),
			expected: false,
		},
		{
			name:     "exactly at threshold",
			now:      time.Date(2026, 3, 15, 16, 55, 0, 0, jst),
			expected: true,
		},
		{
			name:     "within threshold",
			now:      time.Date(2026, 3, 15, 16, 59, 0, 0, jst),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, p.ShouldExtend(tt.now, threshold))
		})
	}
}

func TestAuctionPeriod_Extend(t *testing.T) {
	jst := NewTimeZone(LocationJST).Location()
	end := time.Date(2026, 3, 15, 17, 0, 0, 0, jst)
	p := NewAuctionPeriod(nil, &end)

	duration := 5 * time.Minute
	extendedP := p.Extend(duration)

	expectedEnd := time.Date(2026, 3, 15, 17, 5, 0, 0, jst)
	assert.Equal(t, &expectedEnd, extendedP.EndAt)
}
