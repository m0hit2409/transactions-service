// Package service holds the business logic. Services depend only on the
// repository ports defined in the repository package, never on a concrete
// database or the HTTP layer, which keeps the rules unit-testable in isolation.
package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/repository"
)

type accountService struct {
	accounts repository.AccountRepository
}

func NewAccountService(accounts repository.AccountRepository) AccountService {
	return &accountService{accounts: accounts}
}

// Create validates the document and persists a new account.
func (s *accountService) Create(ctx context.Context, documentNumber string) (domain.Account, error) {
	doc := strings.TrimSpace(documentNumber)
	if doc == "" {
		return domain.Account{}, domain.ErrInvalidDocument
	}

	acc, err := s.accounts.Create(ctx, doc)
	if err != nil {
		return domain.Account{}, fmt.Errorf("create account: %w", err)
	}
	return acc, nil
}

// GetByID returns the account for the given id.
func (s *accountService) GetByID(ctx context.Context, id int64) (domain.Account, error) {
	acc, err := s.accounts.FindByID(ctx, id)
	if err != nil {
		return domain.Account{}, fmt.Errorf("get account: %w", err)
	}
	return acc, nil
}
