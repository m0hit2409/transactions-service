package request

// CreateAccount is the request body for POST /accounts.
type CreateAccount struct {
	DocumentNumber string `json:"document_number" example:"12345678900"`
}
