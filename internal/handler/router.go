package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// NewRouter wires the routes and cross-cutting middleware. New concerns (auth,
// rate limiting, idempotency) are added here as middleware rather than inside
// handlers.
func NewRouter(accounts *AccountHandler, txns *TransactionHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", accounts.Create)
		r.Get("/{accountId}", accounts.GetByID)
	})

	r.Post("/transactions", txns.Create)

	// Interactive API docs.
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("/docs/doc.json")))

	return r
}
