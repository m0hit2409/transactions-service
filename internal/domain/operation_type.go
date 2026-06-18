package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// OperationType is a value object describing a kind of transaction. The Sign
// field (+1 or -1) drives whether the stored amount is a credit or a debit, so
// adding a new operation type is a data change (a new row) rather than a code
// change. IsActive supports soft deletes: a type is deactivated rather than
// removed so historical transactions keep referencing it.
type OperationType struct {
	ID          int       `db:"operation_type_id" json:"operation_type_id"`
	Description string    `db:"description" json:"description"`
	Sign        int       `db:"sign" json:"-"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// Apply returns the amount with this operation type's sign applied. The input is
// always treated as a magnitude (its absolute value), so the recorded sign is
// determined solely by the operation type, never by the caller.
func (o OperationType) Apply(amount decimal.Decimal) decimal.Decimal {
	return amount.Abs().Mul(decimal.NewFromInt(int64(o.Sign)))
}
