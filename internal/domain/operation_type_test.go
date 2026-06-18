package domain_test

import (
	"testing"

	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOperationType_Apply(t *testing.T) {
	tests := []struct {
		name   string
		sign   int
		amount string
		want   string
	}{
		{"debit negates a positive magnitude", -1, "123.45", "-123.45"},
		{"credit keeps a positive magnitude", 1, "60", "60"},
		{"debit normalises a negative input", -1, "-50", "-50"},
		{"credit normalises a negative input", 1, "-50", "50"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ot := domain.OperationType{Sign: tc.sign}
			got := ot.Apply(decimal.RequireFromString(tc.amount))
			assert.True(t, got.Equal(decimal.RequireFromString(tc.want)),
				"want %s, got %s", tc.want, got)
		})
	}
}
