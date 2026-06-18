package validator_test

import (
	"testing"

	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/validator"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func input(amount string) validator.Input {
	return validator.Input{
		Account: domain.Account{ID: 1},
		Amount:  decimal.RequireFromString(amount),
	}
}

func TestPositiveAmountValidator(t *testing.T) {
	v := validator.PositiveAmount{}

	t.Run("accepts a positive amount", func(t *testing.T) {
		assert.NoError(t, v.Validate(input("123.45")))
	})

	t.Run("rejects zero", func(t *testing.T) {
		assert.ErrorIs(t, v.Validate(input("0")), domain.ErrNonPositiveAmount)
	})

	t.Run("rejects a negative amount", func(t *testing.T) {
		assert.ErrorIs(t, v.Validate(input("-10")), domain.ErrNonPositiveAmount)
	})
}

func TestRegistry_ValidatesUsingAllConfiguredRules(t *testing.T) {
	reg := validator.NewRegistry()

	// Every known operation type must reject a non-positive amount.
	for opType := 1; opType <= 4; opType++ {
		in := input("-1")
		in.OperationType = domain.OperationType{ID: opType}
		err := reg.Validate(opType, in)
		require.Error(t, err, "operation type %d should reject negative amount", opType)
		assert.ErrorIs(t, err, domain.ErrNonPositiveAmount)
	}
}

func TestRegistry_PassesValidInput(t *testing.T) {
	reg := validator.NewRegistry()
	in := input("50")
	in.OperationType = domain.OperationType{ID: 1}
	assert.NoError(t, reg.Validate(1, in))
}
