package message

// PushNotificationMessage is the wire format for push notification messages.
// It is used as a data contract between queue producers and worker handlers.
type PushNotificationMessage struct {
	BuyerID int `json:"buyer_id"`
	Payload any `json:"payload"`
}
