package dto

// Invoice DTOs
// InvoiceResponse represents the response body for an invoice.
type InvoiceResponse struct {
	BuyerID     int    `json:"buyer_id"`
	BuyerName   string `json:"buyer_name"`
	TotalAmount int    `json:"total_amount"`
}
