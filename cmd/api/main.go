// Command api is the entrypoint for the transactions service.
//
//	@title			Transactions Service API
//	@version		1.0
//	@description	A small service that manages cardholder accounts and their transactions.
//	@description	Purchases and withdrawals are stored as negative amounts; credit vouchers as positive. The caller always sends a positive amount and the server applies the sign based on the operation type.
//	@BasePath		/
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/m0hit2409/transactions-service/docs" // generated swagger docs
	"github.com/m0hit2409/transactions-service/internal/config"
	"github.com/m0hit2409/transactions-service/internal/handler"
	"github.com/m0hit2409/transactions-service/internal/repository"
	"github.com/m0hit2409/transactions-service/internal/service"
	"github.com/m0hit2409/transactions-service/internal/validator"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	db, err := repository.Connect(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("close db", "error", err)
		}
	}()

	// Wire dependencies top-down: repositories -> services -> handlers -> router.
	accountRepo := repository.NewAccountRepository(db)
	opTypeRepo := repository.NewOperationTypeRepository(db)
	txnRepo := repository.NewTransactionRepository(db)

	accountSvc := service.NewAccountService(accountRepo)
	txnSvc := service.NewTransactionService(accountRepo, opTypeRepo, txnRepo, validator.NewRegistry())

	router := handler.NewRouter(
		handler.NewAccountHandler(accountSvc),
		handler.NewTransactionHandler(txnSvc),
	)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return serve(srv)
}

// serve runs the HTTP server and shuts it down gracefully on SIGINT/SIGTERM.
func serve(srv *http.Server) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
