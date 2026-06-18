package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/service"
	"github.com/m0hit2409/transactions-service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAccountService_Create(t *testing.T) {
	t.Run("creates an account with a valid document", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockAccountRepository(ctrl)
		repo.EXPECT().
			Create(gomock.Any(), "12345678900").
			Return(domain.Account{ID: 1, DocumentNumber: "12345678900"}, nil)

		svc := service.NewAccountService(repo)
		acc, err := svc.Create(context.Background(), "12345678900")

		require.NoError(t, err)
		assert.Equal(t, int64(1), acc.ID)
		assert.Equal(t, "12345678900", acc.DocumentNumber)
	})

	t.Run("rejects an empty document without hitting the repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockAccountRepository(ctrl)
		// No repo call expected.

		svc := service.NewAccountService(repo)
		_, err := svc.Create(context.Background(), "   ")

		assert.ErrorIs(t, err, domain.ErrInvalidDocument)
	})

	t.Run("propagates a duplicate error from the repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockAccountRepository(ctrl)
		repo.EXPECT().
			Create(gomock.Any(), "12345678900").
			Return(domain.Account{}, domain.ErrDuplicateAccount)

		svc := service.NewAccountService(repo)
		_, err := svc.Create(context.Background(), "12345678900")

		assert.ErrorIs(t, err, domain.ErrDuplicateAccount)
	})
}

func TestAccountService_GetByID(t *testing.T) {
	t.Run("returns the account", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockAccountRepository(ctrl)
		repo.EXPECT().
			FindByID(gomock.Any(), int64(1)).
			Return(domain.Account{ID: 1, DocumentNumber: "12345678900"}, nil)

		svc := service.NewAccountService(repo)
		acc, err := svc.GetByID(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, int64(1), acc.ID)
	})

	t.Run("propagates not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockAccountRepository(ctrl)
		repo.EXPECT().
			FindByID(gomock.Any(), int64(99)).
			Return(domain.Account{}, domain.ErrAccountNotFound)

		svc := service.NewAccountService(repo)
		_, err := svc.GetByID(context.Background(), 99)

		assert.ErrorIs(t, err, domain.ErrAccountNotFound)
	})

	t.Run("wraps unexpected repo errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockAccountRepository(ctrl)
		boom := errors.New("connection reset")
		repo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(domain.Account{}, boom)

		svc := service.NewAccountService(repo)
		_, err := svc.GetByID(context.Background(), 1)

		assert.ErrorIs(t, err, boom)
	})
}
