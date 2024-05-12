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

	"github.com/adrianuf22/back-test-psmo/internal/adapter/handler"
	"github.com/adrianuf22/back-test-psmo/internal/adapter/postgres"
	"github.com/adrianuf22/back-test-psmo/internal/config"
	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/adrianuf22/back-test-psmo/internal/domain/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

type App struct {
	cfg        *config.Config
	router     *chi.Mux
	db         *sql.DB
	httpServer *http.Server
	logger     *slog.Logger
}

func main() {
	ctx := context.Background()
	cfg := config.New()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)

	app := App{
		cfg:    cfg,
		router: router,
		httpServer: &http.Server{
			Addr:              ":" + cfg.Http.Port,
			Handler:           router,
			ReadHeaderTimeout: time.Duration(cfg.Http.ReadHeaderTimeout),
		},
		logger: initLogger(),
	}

	app.initDatabase(ctx)
	app.initError(ctx)
	app.initHealth(ctx)
	app.initAccount(ctx)

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

	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Database.TimeoutMS)*time.Second)
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

func (a *App) initAccount(ctx context.Context) {
	repository := postgres.NewAccount(a.db)
	usecase := account.NewUsecase(repository)
	handler.RegisterAccountHandler(ctx, a.router, *usecase)
}

func (a *App) initError(ctx context.Context) {
	handler.RegisterErrorHandler(ctx, a.router)
}

func (a *App) initHealth(ctx context.Context) {
	repository := postgres.NewHealth(a.db)
	usecase := health.NewUsecase(repository)
	handler.RegisterHealthHandler(ctx, a.router, *usecase)
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
