package response

// Me represents the currently authenticated buyer's information.
type Me struct {
	Authenticated bool   `json:"authenticated"`
	BuyerID       int    `json:"buyer_id"`
	Name          string `json:"name"`
}
