package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/models/request"
	"github.com/m0hit2409/transactions-service/internal/models/response"
	"github.com/m0hit2409/transactions-service/internal/service"
)

// AccountHandler exposes the account endpoints.
type AccountHandler struct {
	svc service.AccountService
}

func NewAccountHandler(svc service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

// Create godoc
//
//	@Summary		Create an account
//	@Description	Creates a cardholder account for the given document number.
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateAccount	true	"Account to create"
//	@Success		201		{object}	domain.Account
//	@Failure		400		{object}	response.Error	"malformed body or empty document"
//	@Failure		409		{object}	response.Error	"document number already exists"
//	@Router			/accounts [post]
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateAccount
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, domain.ErrInvalidDocument)
		return
	}

	acc, err := h.svc.Create(r.Context(), req.DocumentNumber)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, acc)
}

// GetByID godoc
//
//	@Summary		Get an account
//	@Description	Retrieves an account by its id.
//	@Tags			accounts
//	@Produce		json
//	@Param			accountId	path		int	true	"Account ID"
//	@Success		200			{object}	domain.Account
//	@Failure		400			{object}	response.Error	"non-numeric id"
//	@Failure		404			{object}	response.Error	"account not found"
//	@Router			/accounts/{accountId} [get]
func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "accountId"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response.Error{Message: "invalid account id"})
		return
	}

	acc, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, acc)
}
