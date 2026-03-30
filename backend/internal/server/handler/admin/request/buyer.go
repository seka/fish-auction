package request

// CreateBuyer holds data for buyer creation.
type CreateBuyer struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
	ContactInfo  string `json:"contact_info"`
}
