package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/handler"
	"github.com/m0hit2409/transactions-service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// newServer builds the full router backed by mock services, so tests exercise
// real routing and middleware, not just a handler func in isolation.
func newServer(t *testing.T) (*httptest.Server, *mocks.MockAccountService, *mocks.MockTransactionService) {
	ctrl := gomock.NewController(t)
	accSvc := mocks.NewMockAccountService(ctrl)
	txSvc := mocks.NewMockTransactionService(ctrl)
	r := handler.NewRouter(handler.NewAccountHandler(accSvc), handler.NewTransactionHandler(txSvc))
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv, accSvc, txSvc
}

func TestCreateAccount(t *testing.T) {
	t.Run("201 with the created account", func(t *testing.T) {
		srv, accSvc, _ := newServer(t)
		accSvc.EXPECT().Create(gomock.Any(), "12345678900").
			Return(domain.Account{ID: 1, DocumentNumber: "12345678900"}, nil)

		resp := post(t, srv.URL+"/accounts", `{"document_number":"12345678900"}`)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var body map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.EqualValues(t, 1, body["account_id"])
		assert.Equal(t, "12345678900", body["document_number"])
	})

	t.Run("409 when the document already exists", func(t *testing.T) {
		srv, accSvc, _ := newServer(t)
		accSvc.EXPECT().Create(gomock.Any(), "12345678900").
			Return(domain.Account{}, domain.ErrDuplicateAccount)

		resp := post(t, srv.URL+"/accounts", `{"document_number":"12345678900"}`)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("400 on a malformed body", func(t *testing.T) {
		srv, _, _ := newServer(t)
		resp := post(t, srv.URL+"/accounts", `{not json`)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetAccount(t *testing.T) {
	t.Run("200 with the account", func(t *testing.T) {
		srv, accSvc, _ := newServer(t)
		accSvc.EXPECT().GetByID(gomock.Any(), int64(1)).
			Return(domain.Account{ID: 1, DocumentNumber: "12345678900"}, nil)

		resp, err := http.Get(srv.URL + "/accounts/1")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("404 when missing", func(t *testing.T) {
		srv, accSvc, _ := newServer(t)
		accSvc.EXPECT().GetByID(gomock.Any(), int64(99)).
			Return(domain.Account{}, domain.ErrAccountNotFound)

		resp, err := http.Get(srv.URL + "/accounts/99")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("400 on a non-numeric id", func(t *testing.T) {
		srv, _, _ := newServer(t)
		resp, err := http.Get(srv.URL + "/accounts/abc")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func post(t *testing.T, url, body string) *http.Response {
	t.Helper()
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	require.NoError(t, err)
	return resp
}
