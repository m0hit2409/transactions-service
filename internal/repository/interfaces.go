// Package repository defines the persistence ports for the domain aggregates and
// holds the PostgreSQL adapters that implement them. The interfaces live here,
// next to their implementation, so the service layer depends on a port rather
// than a concrete database.
package repository

import (
	"context"

	"github.com/m0hit2409/transactions-service/internal/domain"
)

//go:generate mockgen -destination=../../mocks/mock_account_repository.go -package=mocks github.com/m0hit2409/transactions-service/internal/repository AccountRepository

// AccountRepository is the persistence port for accounts.
type AccountRepository interface {
	// Create persists a new account and returns it with its generated ID and
	// timestamps. It returns domain.ErrDuplicateAccount if the document_number exists.
	Create(ctx context.Context, documentNumber string) (domain.Account, error)

	// FindByID returns the account or domain.ErrAccountNotFound.
	FindByID(ctx context.Context, id int64) (domain.Account, error)
}

//go:generate mockgen -destination=../../mocks/mock_transaction_repository.go -package=mocks github.com/m0hit2409/transactions-service/internal/repository TransactionRepository

// TransactionRepository is the persistence port for transactions.
type TransactionRepository interface {
	// Create persists the transaction, sets the server-side event_date, and
	// returns it with its generated ID and timestamp.
	Create(ctx context.Context, tx domain.Transaction) (domain.Transaction, error)
}

//go:generate mockgen -destination=../../mocks/mock_operation_type_repository.go -package=mocks github.com/m0hit2409/transactions-service/internal/repository OperationTypeRepository

// OperationTypeRepository is the persistence port for operation types. These are
// static seed data, so there is no Create — only lookup.
type OperationTypeRepository interface {
	// FindByID returns the operation type or domain.ErrOperationTypeNotFound.
	FindByID(ctx context.Context, id int) (domain.OperationType, error)
}
