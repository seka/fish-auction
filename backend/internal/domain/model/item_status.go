package model

// ItemStatus represents the status of an auction item
type ItemStatus string

const (
	// ItemStatusPending provides ItemStatusPending related functionality.
	ItemStatusPending   ItemStatus = "Pending"
	// ItemStatusAvailable provides ItemStatusAvailable related functionality.
	ItemStatusAvailable ItemStatus = "Available"
	// ItemStatusSold provides ItemStatusSold related functionality.
	ItemStatusSold      ItemStatus = "Sold"
)

// IsValid checks if the status is valid
func (s ItemStatus) IsValid() bool {
	switch s {
	case ItemStatusPending, ItemStatusAvailable, ItemStatusSold:
		return true
	}
	return false
}

// String returns the string representation
func (s ItemStatus) String() string {
	return string(s)
}
