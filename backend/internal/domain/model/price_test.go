package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBidPrice_CalculateMinIncrement(t *testing.T) {
	tests := []struct {
		name     string
		price    int
		expected int
	}{
		{"Under 1,000", 999, 100},
		{"Exactly 1,000", 1000, 500},
		{"Exactly 9,999", 9999, 500},
		{"Exactly 10,000", 10000, 1000},
		{"Exactly 99,999", 99999, 1000},
		{"Exactly 100,000", 100000, 5000},
		{"Over 100,000", 150000, 5000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewBidPrice(tt.price)
			assert.Equal(t, tt.expected, p.CalculateMinIncrement().Value)
		})
	}
}

func TestBidPrice_Add(t *testing.T) {
	p1 := NewBidPrice(1000)
	p2 := NewBidPrice(500)
	result := p1.Add(p2)
	assert.Equal(t, 1500, result.Value)
}

func TestBidPrice_LessThan(t *testing.T) {
	p1 := NewBidPrice(1000)
	p2 := NewBidPrice(2000)
	assert.True(t, p1.LessThan(p2))
	assert.False(t, p2.LessThan(p1))
	assert.False(t, p1.LessThan(p1))
}
