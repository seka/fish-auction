package model

const (
	// MinBidIncrementUnder1k provides MinBidIncrementUnder1k related functionality.
	MinBidIncrementUnder1k = 100
	// MinBidIncrementUnder10k provides MinBidIncrementUnder10k related functionality.
	MinBidIncrementUnder10k = 500
	// MinBidIncrementUnder100k provides MinBidIncrementUnder100k related functionality.
	MinBidIncrementUnder100k = 1000
	// MinBidIncrementDefault represents a minbidincrementdefault in the system.
	MinBidIncrementDefault = 5000
)

// BidPrice represents a price of a bid.
type BidPrice struct {
	Value int
}

// NewBidPrice creates a new BidPrice.
func NewBidPrice(amount int) BidPrice {
	return BidPrice{Value: amount}
}

// Amount returns the integer value.
func (p BidPrice) Amount() int {
	return p.Value
}

// CalculateMinIncrement returns the minimum increment allowed for this price.
func (p BidPrice) CalculateMinIncrement() BidPrice {
	if p.Value < 1000 {
		return NewBidPrice(MinBidIncrementUnder1k)
	}
	if p.Value < 10000 {
		return NewBidPrice(MinBidIncrementUnder10k)
	}
	if p.Value < 100000 {
		return NewBidPrice(MinBidIncrementUnder100k)
	}
	return NewBidPrice(MinBidIncrementDefault)
}

// Add returns a new BidPrice with the added amount.
func (p BidPrice) Add(other BidPrice) BidPrice {
	return NewBidPrice(p.Value + other.Value)
}

// LessThan returns true if p is less than other.
func (p BidPrice) LessThan(other BidPrice) bool {
	return p.Value < other.Value
}
