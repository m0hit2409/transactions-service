package domain

import "errors"

// Sentinel errors expressed in domain terms. The HTTP layer is responsible for
// translating these into status codes (see internal/handler/response.go), so the
// domain and service layers never import net/http.
var (
	// ErrAccountNotFound is returned when an account lookup yields no row.
	ErrAccountNotFound = errors.New("account not found")

	// ErrDuplicateAccount is returned when creating an account whose
	// document_number already exists.
	ErrDuplicateAccount = errors.New("account already exists")

	// ErrInvalidDocument is returned when a document_number is empty or malformed.
	ErrInvalidDocument = errors.New("invalid document number")

	// ErrOperationTypeNotFound is returned when a transaction references an
	// operation_type_id that does not exist.
	ErrOperationTypeNotFound = errors.New("invalid operation type")

	// ErrNonPositiveAmount is returned when a transaction amount is zero or negative.
	// The caller must always send a positive amount; the service owns the sign.
	ErrNonPositiveAmount = errors.New("amount must be positive")
)
