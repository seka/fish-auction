package model

import (
	"fmt"
)

// JobType represents the type of an asynchronous job.
type JobType string

const (
	// JobTypePushNotification is the job type for sending push notifications.
	JobTypePushNotification JobType = "push_notification"
)

// NewJobType creates a JobType from a string and validates it.
func NewJobType(s string) (JobType, error) {
	switch JobType(s) {
	case JobTypePushNotification:
		return JobTypePushNotification, nil
	default:
		return "", fmt.Errorf("unsupported job type: %s", s)
	}
}
