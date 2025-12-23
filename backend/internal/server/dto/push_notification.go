package dto

// SubscribePushRequest represents the request body for subscribing to push notifications
type SubscribePushRequest struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}
