// Package validator holds the composable rules applied to a transaction before
// it is persisted. Rules are small and independently testable; each operation
// type declares which rules apply to it via the Registry. This models the real
// variation (rules differ per type) without one near-identical Strategy class
// per operation type.
package validator

import (
	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/shopspring/decimal"
)

// Input is everything a rule may need to make a decision. Amount is the raw
// magnitude supplied by the caller (sign not yet applied).
type Input struct {
	Account       domain.Account
	OperationType domain.OperationType
	Amount        decimal.Decimal
}

// TransactionValidator is a single rule.
type TransactionValidator interface {
	Validate(in Input) error
}
