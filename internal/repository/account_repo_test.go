package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return sqlx.NewDb(db, "sqlmock"), mock
}

func TestAccountRepo_Create(t *testing.T) {
	cols := []string{"account_id", "document_number", "is_active", "created_at"}

	t.Run("success", func(t *testing.T) {
		db, mock := newMockDB(t)
		now := time.Now()

		mock.ExpectQuery(`INSERT INTO accounts`).
			WithArgs("12345678900").
			WillReturnRows(sqlmock.NewRows(cols).AddRow(1, "12345678900", true, now))

		acc, err := repository.NewAccountRepository(db).Create(context.Background(), "12345678900")

		require.NoError(t, err)
		assert.Equal(t, int64(1), acc.ID)
		assert.Equal(t, "12345678900", acc.DocumentNumber)
		assert.True(t, acc.IsActive)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate document returns ErrDuplicateAccount", func(t *testing.T) {
		db, mock := newMockDB(t)

		mock.ExpectQuery(`INSERT INTO accounts`).
			WithArgs("12345678900").
			WillReturnError(sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique})

		_, err := repository.NewAccountRepository(db).Create(context.Background(), "12345678900")

		assert.ErrorIs(t, err, domain.ErrDuplicateAccount)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error is wrapped", func(t *testing.T) {
		db, mock := newMockDB(t)
		dbErr := errors.New("connection lost")

		mock.ExpectQuery(`INSERT INTO accounts`).
			WithArgs("12345678900").
			WillReturnError(dbErr)

		_, err := repository.NewAccountRepository(db).Create(context.Background(), "12345678900")

		assert.ErrorContains(t, err, "insert account")
		assert.ErrorIs(t, err, dbErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAccountRepo_FindByID(t *testing.T) {
	cols := []string{"account_id", "document_number", "is_active", "created_at"}

	t.Run("success", func(t *testing.T) {
		db, mock := newMockDB(t)
		now := time.Now()

		mock.ExpectQuery(`SELECT account_id`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows(cols).AddRow(1, "12345678900", true, now))

		acc, err := repository.NewAccountRepository(db).FindByID(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, int64(1), acc.ID)
		assert.Equal(t, "12345678900", acc.DocumentNumber)
		assert.True(t, acc.IsActive)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found returns ErrAccountNotFound", func(t *testing.T) {
		db, mock := newMockDB(t)

		mock.ExpectQuery(`SELECT account_id`).
			WithArgs(int64(99)).
			WillReturnRows(sqlmock.NewRows(cols))

		_, err := repository.NewAccountRepository(db).FindByID(context.Background(), 99)

		assert.ErrorIs(t, err, domain.ErrAccountNotFound)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error is wrapped", func(t *testing.T) {
		db, mock := newMockDB(t)
		dbErr := errors.New("connection lost")

		mock.ExpectQuery(`SELECT account_id`).
			WithArgs(int64(1)).
			WillReturnError(dbErr)

		_, err := repository.NewAccountRepository(db).FindByID(context.Background(), 1)

		assert.ErrorContains(t, err, "select account")
		assert.ErrorIs(t, err, dbErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
