package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// CreateTransactionInput carries the raw inputs for creating a transaction.
type CreateTransactionInput struct {
	AccountID       int64
	OperationTypeID int
	Amount          decimal.Decimal
}

// Transaction is an immutable financial fact associated with an account. Amount
// is stored already-signed: the posted value is a snapshot of the truth at write
// time and does not depend on the (mutable) operation_types.sign being read back.
type Transaction struct {
	ID              int64           `db:"transaction_id" json:"transaction_id"`
	AccountID       int64           `db:"account_id" json:"account_id"`
	OperationTypeID int             `db:"operation_type_id" json:"operation_type_id"`
	Amount          decimal.Decimal `db:"amount" json:"amount"`
	EventDate       time.Time       `db:"event_date" json:"event_date"`
}

// NewTransaction builds a Transaction whose amount has the operation type's sign
// applied. EventDate is intentionally left zero here and set by the repository
// at insert time so the server owns the event timestamp.
func NewTransaction(account Account, opType OperationType, amount decimal.Decimal) Transaction {
	return Transaction{
		AccountID:       account.ID,
		OperationTypeID: opType.ID,
		Amount:          opType.Apply(amount),
	}
}
