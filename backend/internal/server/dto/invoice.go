package dto

// InvoiceResponse is a data transfer object.
type InvoiceResponse struct {
	BuyerID     int    `json:"buyer_id"`
	BuyerName   string `json:"buyer_name"`
	TotalAmount int    `json:"total_amount"`
}
