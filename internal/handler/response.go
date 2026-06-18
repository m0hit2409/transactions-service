package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/m0hit2409/transactions-service/internal/domain"
	"github.com/m0hit2409/transactions-service/internal/models/response"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response", "error", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	status, msg := statusFor(err)
	if status >= http.StatusInternalServerError {
		slog.Error("request failed", "error", err)
		msg = "internal server error"
	}
	writeJSON(w, status, response.Error{Message: msg})
}

func statusFor(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		return http.StatusNotFound, domain.ErrAccountNotFound.Error()
	case errors.Is(err, domain.ErrDuplicateAccount):
		return http.StatusConflict, domain.ErrDuplicateAccount.Error()
	case errors.Is(err, domain.ErrInvalidDocument):
		return http.StatusBadRequest, domain.ErrInvalidDocument.Error()
	case errors.Is(err, domain.ErrNonPositiveAmount):
		return http.StatusBadRequest, domain.ErrNonPositiveAmount.Error()
	case errors.Is(err, domain.ErrOperationTypeNotFound):
		return http.StatusUnprocessableEntity, domain.ErrOperationTypeNotFound.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}
