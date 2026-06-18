package domain

import (
	"time"
)

// Account is the cardholder's account aggregate. IsActive supports soft deletes:
// an account is deactivated rather than removed so its transaction history stays
// intact.
type Account struct {
	ID             int64     `db:"account_id" json:"account_id"`
	DocumentNumber string    `db:"document_number" json:"document_number"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
