package util

import "time"

// FormatTimestamp formats a time.Time pointer into an RFC3339 string pointer.
func FormatTimestamp(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
