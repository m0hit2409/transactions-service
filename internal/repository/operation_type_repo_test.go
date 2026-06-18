package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationTypeRepo_FindByID(t *testing.T) {
	cols := []string{"operation_type_id", "description", "sign", "is_active", "created_at"}

	t.Run("success", func(t *testing.T) {
		db, mock := newMockDB(t)
		now := time.Now()

		mock.ExpectQuery(`SELECT operation_type_id`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows(cols).AddRow(1, "Normal Purchase", -1, true, now))

		ot, err := repository.NewOperationTypeRepository(db).FindByID(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, 1, ot.ID)
		assert.Equal(t, "Normal Purchase", ot.Description)
		assert.Equal(t, -1, ot.Sign)
		assert.True(t, ot.IsActive)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found returns ErrOperationTypeNotFound", func(t *testing.T) {
		db, mock := newMockDB(t)

		mock.ExpectQuery(`SELECT operation_type_id`).
			WithArgs(int64(99)).
			WillReturnRows(sqlmock.NewRows(cols))

		_, err := repository.NewOperationTypeRepository(db).FindByID(context.Background(), 99)

		assert.ErrorIs(t, err, domain.ErrOperationTypeNotFound)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error is wrapped", func(t *testing.T) {
		db, mock := newMockDB(t)
		dbErr := errors.New("connection lost")

		mock.ExpectQuery(`SELECT operation_type_id`).
			WithArgs(int64(99)).
			WillReturnError(dbErr)

		_, err := repository.NewOperationTypeRepository(db).FindByID(context.Background(), 99)

		assert.ErrorContains(t, err, "select operation type")
		assert.ErrorIs(t, err, dbErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
