package handler

import (
	"encoding/json"
	"net/http"

	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/models/request"
	"github.com/m0hit2409/transactions-service/internal/models/response"
	"github.com/m0hit2409/transactions-service/internal/service"
)

// TransactionHandler exposes the transaction endpoint.
type TransactionHandler struct {
	svc service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// Create godoc
//
//	@Summary		Create a transaction
//	@Description	Records a transaction against an account. The caller always sends a positive amount; the server applies the sign based on the operation type (purchases and withdrawals are stored negative, credit vouchers positive).
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateTransaction	true	"Transaction to create"
//	@Success		201		{object}	domain.Transaction
//	@Failure		400		{object}	response.Error	"malformed body or non-positive amount"
//	@Failure		404		{object}	response.Error	"account not found"
//	@Failure		422		{object}	response.Error	"invalid operation type"
//	@Router			/transactions [post]
func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateTransaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, response.Error{Message: "invalid request body"})
		return
	}

	tx, err := h.svc.Create(r.Context(), domain.CreateTransactionInput{
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Amount:          req.Amount,
	})

	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, tx)
}
