package repository

// The adapters below implement the persistence ports against SQLite using
// sqlx and raw SQL. Database errors are translated into domain sentinel errors
// here so nothing above this layer needs to know about SQL state codes.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/m0hit2409/transactions-service/internal/domain"
)

type accountRepo struct {
	db *sqlx.DB
}

// NewAccountRepository returns a SQLite-backed AccountRepository.
func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) Create(ctx context.Context, documentNumber string) (domain.Account, error) {
	const q = `
		INSERT INTO accounts (document_number)
		VALUES (?)
		RETURNING account_id, document_number, is_active, created_at`

	var acc domain.Account
	if err := r.db.GetContext(ctx, &acc, q, documentNumber); err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return domain.Account{}, domain.ErrDuplicateAccount
		}
		return domain.Account{}, fmt.Errorf("insert account: %w", err)
	}
	return acc, nil
}

func (r *accountRepo) FindByID(ctx context.Context, id int64) (domain.Account, error) {
	const q = `
		SELECT account_id, document_number, is_active, created_at
		FROM accounts
		WHERE account_id = ?`

	var acc domain.Account
	if err := r.db.GetContext(ctx, &acc, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Account{}, domain.ErrAccountNotFound
		}
		return domain.Account{}, fmt.Errorf("select account: %w", err)
	}
	return acc, nil
}
