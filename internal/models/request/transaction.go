package request

import "github.com/shopspring/decimal"

// CreateTransaction is the request body for POST /transactions.
type CreateTransaction struct {
	AccountID       int64           `json:"account_id" example:"1"`
	OperationTypeID int             `json:"operation_type_id" example:"4"`
	Amount          decimal.Decimal `json:"amount" swaggertype:"number" example:"123.45"`
}
