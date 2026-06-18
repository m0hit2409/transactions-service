package validator

import (
	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/shopspring/decimal"
)

// PositiveAmount rejects amounts that are zero or negative. The API contract is
// that callers always send a positive magnitude; the service owns the sign.
type PositiveAmount struct{}

func (PositiveAmount) Validate(in Input) error {
	if in.Amount.LessThanOrEqual(decimal.Zero) {
		return domain.ErrNonPositiveAmount
	}
	return nil
}
