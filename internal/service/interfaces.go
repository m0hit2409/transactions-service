//go:generate mockgen -source=interfaces.go -destination=../../mocks/mock_services.go -package=mocks

package service

import (
	"context"

	"github.com/m0hit2409/transactions-service/internal/domain"
)

// AccountService is the contract for account use cases.
type AccountService interface {
	Create(ctx context.Context, documentNumber string) (domain.Account, error)
	GetByID(ctx context.Context, id int64) (domain.Account, error)
}

// TransactionService is the contract for transaction use cases.
type TransactionService interface {
	Create(ctx context.Context, cmd domain.CreateTransactionInput) (domain.Transaction, error)
}
