package model

import "time"

// LocationName represents a supported time zone location name.
type LocationName string

const (
	LocationJST LocationName = "Asia/Tokyo"
)

var (
	offsets = map[LocationName]int{
		LocationJST: 9 * 60 * 60,
	}
)

// TimeZone represents a time zone used for auction operations.
type TimeZone struct {
	location *time.Location
}

// NewTimeZone returns a TimeZone for the given location name.
func NewTimeZone(name LocationName) TimeZone {
	offset, ok := offsets[name]
	if !ok {
		// Default to JST offset if name is unknown.
		offset = offsets[LocationJST]
		name = LocationJST
	}

	return TimeZone{
		location: time.FixedZone(string(name), offset),
	}
}

// Now returns the current time in this TimeZone.
func (tz TimeZone) Now() time.Time {
	return time.Now().In(tz.location)
}

// At returns the given time in this TimeZone.
func (tz TimeZone) At(t time.Time) time.Time {
	return t.In(tz.location)
}

// Location returns the underlying time.Location.
func (tz TimeZone) Location() *time.Location {
	return tz.location
}
