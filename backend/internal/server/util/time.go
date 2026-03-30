package util

import "time"

// FormatTime formats a time.Time pointer into a string pointer ("15:04:05").
func FormatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("15:04:05")
	return &s
}
