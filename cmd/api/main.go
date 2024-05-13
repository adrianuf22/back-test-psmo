package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adrianuf22/back-test-psmo/internal/adapter/mux"
	"github.com/adrianuf22/back-test-psmo/internal/adapter/postgres"
	"github.com/adrianuf22/back-test-psmo/internal/config"
	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/adrianuf22/back-test-psmo/internal/domain/health"
	"github.com/adrianuf22/back-test-psmo/internal/domain/transaction"
	_ "github.com/lib/pq"
)

type App struct {
	cfg        *config.Config
	router     *http.ServeMux
	db         *sql.DB
	httpServer *http.Server
	logger     *slog.Logger
}

func main() {
	ctx := context.Background()
	cfg := config.New()

	router := http.NewServeMux()

	app := App{
		cfg:    cfg,
		router: router,
		httpServer: &http.Server{
			Addr:         ":" + cfg.Http.Port,
			Handler:      router,
			ReadTimeout:  time.Duration(cfg.Http.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(cfg.Http.WriteTimeout) * time.Millisecond,
		},
		logger: initLogger(),
	}

	app.initDatabase(ctx)
	app.initHandlers(ctx)

	app.run(ctx)
}

func (a *App) initDatabase(ctx context.Context) {
	var err error
	a.db, err = sql.Open(a.cfg.Database.Driver, a.cfg.Database.DSN())
	if err != nil {
		panic(fmt.Errorf(`error on opening db connection %w`, err))
	}

	a.db.SetMaxOpenConns(a.cfg.Database.MaxOpenConns)
	a.db.SetMaxIdleConns(a.cfg.Database.MaxIdleConns)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Database.Timeout)*time.Millisecond)
	defer cancel()

	err = a.db.PingContext(ctx)
	if err != nil {
		panic(fmt.Errorf(`error on verifying db connection %w`, err))
	}
}

func initLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			}

			if a.Key == slog.MessageKey {
				a.Key = "message"
			}

			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func (a *App) initHandlers(ctx context.Context) {
	// Repositories
	accountRepository := postgres.NewAccountRepo(a.db)
	transactionRepository := postgres.NewTransactionRepo(a.db)
	healthRepository := postgres.NewHealthRepo(a.db)

	// Usecases
	accountUsecase := account.NewUsecase(accountRepository)
	transactionUsecase := transaction.NewUsecase(transactionRepository, accountRepository)
	healthUsecase := health.NewUsecase(healthRepository)

	mux.RegisterAccountHandler(ctx, a.router, accountUsecase)
	mux.RegisterTransactionHandler(ctx, a.router, transactionUsecase)
	mux.RegisterErrorHandler(ctx, a.router)
	mux.RegisterHealthHandler(ctx, a.router, *healthUsecase)
}

func (a *App) run(ctx context.Context) {
	go func() {
		slog.Info("Running...", "addr", a.httpServer.Addr, "pid", os.Getpid())
		err := a.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic("error on start HTTP server")
		}
	}()

	a.gracefulShutdown(ctx)
}

func (a *App) gracefulShutdown(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	slog.Info("Shutting down...")

	ctx, shutdown := context.WithTimeout(ctx, 3*time.Second)
	defer a.closeResources(shutdown)

	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		panic("error on shutting down")
	}
}

func (a *App) closeResources(shutdown context.CancelFunc) {
	a.db.Close()
	shutdown()
}
