package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionRepo_Create(t *testing.T) {
	cols := []string{"transaction_id", "account_id", "operation_type_id", "amount", "event_date"}
	amount := decimal.NewFromFloat(-100.50)

	t.Run("success", func(t *testing.T) {
		db, mock := newMockDB(t)
		now := time.Now()

		mock.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(int64(1), int64(1), amount.String()).
			WillReturnRows(sqlmock.NewRows(cols).AddRow(1, 1, 1, amount.String(), now))

		tx := domain.Transaction{AccountID: 1, OperationTypeID: 1, Amount: amount}
		created, err := repository.NewTransactionRepository(db).Create(context.Background(), tx)

		require.NoError(t, err)
		assert.Equal(t, int64(1), created.ID)
		assert.Equal(t, int64(1), created.AccountID)
		assert.Equal(t, 1, created.OperationTypeID)
		assert.True(t, amount.Equal(created.Amount))
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error is wrapped", func(t *testing.T) {
		db, mock := newMockDB(t)
		dbErr := errors.New("connection lost")

		mock.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(int64(1), int64(1), amount.String()).
			WillReturnError(dbErr)

		tx := domain.Transaction{AccountID: 1, OperationTypeID: 1, Amount: amount}
		_, err := repository.NewTransactionRepository(db).Create(context.Background(), tx)

		assert.ErrorContains(t, err, "insert transaction")
		assert.ErrorIs(t, err, dbErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
