package model

// JobType represents the type of an asynchronous job.
type JobType string

const (
	// JobTypePushNotification is the job type for sending push notifications.
	JobTypePushNotification JobType = "push_notification"
)
