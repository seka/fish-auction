package event

// PushPayload is the data delivered to the browser Service Worker.
// フロントエンド (frontend/public/sw.js) が title / body / url を直接参照するため、
// ここで wire format を固定する。
type PushPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
}

// PushNotificationMessage is the wire format for push notification jobs.
// It is used as a data contract between queue producers and worker handlers
// across all JobTypePush* kinds.
type PushNotificationMessage struct {
	BuyerID int         `json:"buyer_id"`
	Payload PushPayload `json:"payload"`
}
