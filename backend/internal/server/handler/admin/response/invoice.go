package response

// Invoice represents a summary of a buyer's total purchases for admins.
type Invoice struct {
	BuyerID     int    `json:"buyer_id"`
	BuyerName   string `json:"buyer_name"`
	TotalAmount int    `json:"total_amount"`
}
