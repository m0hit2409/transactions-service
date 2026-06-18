package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/service"
	"github.com/m0hit2409/transactions-service/internal/validator"
	"github.com/m0hit2409/transactions-service/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type txDeps struct {
	accounts *mocks.MockAccountRepository
	opTypes  *mocks.MockOperationTypeRepository
	txns     *mocks.MockTransactionRepository
	svc      service.TransactionService
}

func newTxDeps(t *testing.T) txDeps {
	ctrl := gomock.NewController(t)
	accounts := mocks.NewMockAccountRepository(ctrl)
	opTypes := mocks.NewMockOperationTypeRepository(ctrl)
	txns := mocks.NewMockTransactionRepository(ctrl)
	svc := service.NewTransactionService(accounts, opTypes, txns, validator.NewRegistry())
	return txDeps{accounts, opTypes, txns, svc}
}

func cmd(opType int, amount string) domain.CreateTransactionInput {
	return domain.CreateTransactionInput{
		AccountID:       1,
		OperationTypeID: opType,
		Amount:          decimal.RequireFromString(amount),
	}
}

func TestTransactionService_Create(t *testing.T) {
	t.Run("applies negative sign to a purchase", func(t *testing.T) {
		d := newTxDeps(t)
		d.accounts.EXPECT().FindByID(gomock.Any(), int64(1)).Return(domain.Account{ID: 1}, nil)
		d.opTypes.EXPECT().FindByID(gomock.Any(), 1).
			Return(domain.OperationType{ID: 1, Description: "Normal Purchase", Sign: -1}, nil)
		d.txns.EXPECT().Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, tx domain.Transaction) (domain.Transaction, error) {
				assert.True(t, tx.Amount.Equal(decimal.RequireFromString("-123.45")),
					"expected -123.45, got %s", tx.Amount)
				tx.ID = 10
				return tx, nil
			})

		got, err := d.svc.Create(context.Background(), cmd(1, "123.45"))

		require.NoError(t, err)
		assert.Equal(t, int64(10), got.ID)
		assert.True(t, got.Amount.Equal(decimal.RequireFromString("-123.45")))
	})

	t.Run("applies positive sign to a credit voucher", func(t *testing.T) {
		d := newTxDeps(t)
		d.accounts.EXPECT().FindByID(gomock.Any(), int64(1)).Return(domain.Account{ID: 1}, nil)
		d.opTypes.EXPECT().FindByID(gomock.Any(), 4).
			Return(domain.OperationType{ID: 4, Description: "Credit Voucher", Sign: 1}, nil)
		d.txns.EXPECT().Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, tx domain.Transaction) (domain.Transaction, error) {
				assert.True(t, tx.Amount.Equal(decimal.RequireFromString("60")))
				return tx, nil
			})

		_, err := d.svc.Create(context.Background(), cmd(4, "60"))
		require.NoError(t, err)
	})

	t.Run("rejects a non-positive amount", func(t *testing.T) {
		d := newTxDeps(t)
		d.accounts.EXPECT().FindByID(gomock.Any(), int64(1)).Return(domain.Account{ID: 1}, nil)
		d.opTypes.EXPECT().FindByID(gomock.Any(), 1).
			Return(domain.OperationType{ID: 1, Sign: -1}, nil)

		_, err := d.svc.Create(context.Background(), cmd(1, "-10"))
		assert.ErrorIs(t, err, domain.ErrNonPositiveAmount)
	})

	t.Run("propagates account not found", func(t *testing.T) {
		d := newTxDeps(t)
		d.accounts.EXPECT().FindByID(gomock.Any(), int64(1)).
			Return(domain.Account{}, domain.ErrAccountNotFound)

		_, err := d.svc.Create(context.Background(), cmd(1, "50"))
		assert.ErrorIs(t, err, domain.ErrAccountNotFound)
	})

	t.Run("propagates operation type not found", func(t *testing.T) {
		d := newTxDeps(t)
		d.accounts.EXPECT().FindByID(gomock.Any(), int64(1)).Return(domain.Account{ID: 1}, nil)
		d.opTypes.EXPECT().FindByID(gomock.Any(), 99).
			Return(domain.OperationType{}, domain.ErrOperationTypeNotFound)

		_, err := d.svc.Create(context.Background(), cmd(99, "50"))
		assert.ErrorIs(t, err, domain.ErrOperationTypeNotFound)
	})

	t.Run("wraps a repo error from transaction insert", func(t *testing.T) {
		d := newTxDeps(t)
		dbErr := errors.New("connection lost")
		d.accounts.EXPECT().FindByID(gomock.Any(), int64(1)).Return(domain.Account{ID: 1}, nil)
		d.opTypes.EXPECT().FindByID(gomock.Any(), 1).
			Return(domain.OperationType{ID: 1, Sign: -1}, nil)
		d.txns.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.Transaction{}, dbErr)

		_, err := d.svc.Create(context.Background(), cmd(1, "50"))
		assert.ErrorIs(t, err, dbErr)
	})
}
