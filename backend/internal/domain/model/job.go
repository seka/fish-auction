package model

// JobType represents the type of an asynchronous job.
type JobType string

const (
	// JobTypePushNotification is the job type for sending push notifications.
	JobTypePushNotification JobType = "push_notification"
)

// PushNotificationJob represents the domain model for a push notification job.
// Note: This model is used within the domain and usecase layers and does not contain JSON tags.
type PushNotificationJob struct {
	BuyerID int
	Payload any
}
