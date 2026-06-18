package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTransaction(t *testing.T) {
	t.Run("201 with the created transaction", func(t *testing.T) {
		srv, _, txSvc := newServer(t)
		txSvc.EXPECT().
			Create(gomock.Any(), domain.CreateTransactionInput{
				AccountID:       1,
				OperationTypeID: 1,
				Amount:          decimal.RequireFromString("123.45"),
			}).
			Return(domain.Transaction{ID: 10, AccountID: 1, OperationTypeID: 1, Amount: decimal.RequireFromString("-123.45")}, nil)

		resp := post(t, srv.URL+"/transactions",
			`{"account_id":1,"operation_type_id":1,"amount":123.45}`)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var body map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.EqualValues(t, 10, body["transaction_id"])
		assert.Equal(t, "-123.45", body["amount"])
	})

	t.Run("404 when the account is missing", func(t *testing.T) {
		srv, _, txSvc := newServer(t)
		txSvc.EXPECT().Create(gomock.Any(), gomock.Any()).
			Return(domain.Transaction{}, domain.ErrAccountNotFound)

		resp := post(t, srv.URL+"/transactions",
			`{"account_id":99,"operation_type_id":1,"amount":50}`)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("422 on an invalid operation type", func(t *testing.T) {
		srv, _, txSvc := newServer(t)
		txSvc.EXPECT().Create(gomock.Any(), gomock.Any()).
			Return(domain.Transaction{}, domain.ErrOperationTypeNotFound)

		resp := post(t, srv.URL+"/transactions",
			`{"account_id":1,"operation_type_id":99,"amount":50}`)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("400 on a non-positive amount", func(t *testing.T) {
		srv, _, txSvc := newServer(t)
		txSvc.EXPECT().Create(gomock.Any(), gomock.Any()).
			Return(domain.Transaction{}, domain.ErrNonPositiveAmount)

		resp := post(t, srv.URL+"/transactions",
			`{"account_id":1,"operation_type_id":1,"amount":-10}`)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("400 on a malformed body", func(t *testing.T) {
		srv, _, _ := newServer(t)
		resp := post(t, srv.URL+"/transactions", `{bad`)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
