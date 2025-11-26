package model

// ItemStatus represents the status of an auction item
type ItemStatus string

const (
	ItemStatusAvailable ItemStatus = "Available"
	ItemStatusSold      ItemStatus = "Sold"
)

// IsValid checks if the status is valid
func (s ItemStatus) IsValid() bool {
	switch s {
	case ItemStatusAvailable, ItemStatusSold:
		return true
	}
	return false
}

// String returns the string representation
func (s ItemStatus) String() string {
	return string(s)
}
