package service

import (
	"context"
	"fmt"

	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/repository"
	"github.com/m0hit2409/transactions-service/internal/validator"
)

// transactionService implements the transaction use case.
type transactionService struct {
	accounts   repository.AccountRepository
	opTypes    repository.OperationTypeRepository
	txns       repository.TransactionRepository
	validators *validator.Registry
}

func NewTransactionService(
	accounts repository.AccountRepository,
	opTypes repository.OperationTypeRepository,
	txns repository.TransactionRepository,
	validators *validator.Registry,
) TransactionService {
	return &transactionService{
		accounts:   accounts,
		opTypes:    opTypes,
		txns:       txns,
		validators: validators,
	}
}

// Create loads the referenced account and operation type, runs the validation
// rules for that operation type, applies the sign, and persists the transaction.
func (s *transactionService) Create(ctx context.Context, cmd domain.CreateTransactionInput) (domain.Transaction, error) {
	account, err := s.accounts.FindByID(ctx, cmd.AccountID)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("load account: %w", err)
	}

	opType, err := s.opTypes.FindByID(ctx, cmd.OperationTypeID)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("load operation type: %w", err)
	}

	in := validator.Input{Account: account, OperationType: opType, Amount: cmd.Amount}
	if err := s.validators.Validate(opType.ID, in); err != nil {
		return domain.Transaction{}, fmt.Errorf("validate transaction: %w", err)
	}

	created, err := s.txns.Create(ctx, domain.NewTransaction(account, opType, cmd.Amount))
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("create transaction: %w", err)
	}
	return created, nil
}
