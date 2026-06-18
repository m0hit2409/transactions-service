package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/m0hit2409/transactions-service/internal/domain"
)

type transactionRepo struct {
	db *sqlx.DB
}

// NewTransactionRepository returns a SQLite-backed TransactionRepository.
func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, tx domain.Transaction) (domain.Transaction, error) {
	// event_date is intentionally omitted so the column default (now()) sets it:
	// the server owns the event timestamp.
	const q = `
		INSERT INTO transactions (account_id, operation_type_id, amount)
		VALUES (?, ?, ?)
		RETURNING transaction_id, account_id, operation_type_id, amount, event_date`

	var created domain.Transaction
	err := r.db.GetContext(ctx, &created, q,
		tx.AccountID, tx.OperationTypeID, tx.Amount)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("insert transaction: %w", err)
	}
	return created, nil
}
