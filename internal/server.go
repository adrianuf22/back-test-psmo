package internal

import (
	"context"
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
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Cfg        *config.Config
	Router     *http.ServeMux
	Db         *pgxpool.Conn
	HttpServer *http.Server
	Logger     *slog.Logger
}

func NewServer(ctx context.Context) Server {
	cfg := config.New()

	router := http.NewServeMux()

	server := Server{
		Cfg:    cfg,
		Router: router,
		HttpServer: &http.Server{
			Addr:         ":" + cfg.Http.Port,
			Handler:      router,
			ReadTimeout:  time.Duration(cfg.Http.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(cfg.Http.WriteTimeout) * time.Millisecond,
		},
		Logger: slog.Default(),
	}

	server.initDatabase(ctx)
	server.initHandlers(ctx)

	return server
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		slog.Info("Running...", "addr", s.HttpServer.Addr, "pid", os.Getpid())
		err := s.HttpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic("error on start HTTP server")
		}
	}()

	s.gracefulShutdown(ctx)
}

func (s *Server) initDatabase(ctx context.Context) {
	pool, err := pgxpool.New(ctx, s.Cfg.Database.Dsn)
	if err != nil {
		panic(fmt.Errorf(`error on starting pool connection %w`, err))
	}

	s.Db, err = pool.Acquire(ctx)
	if err != nil {
		panic(fmt.Errorf(`error on opening db connection %w`, err))
	}

	// s.db.SetMaxOpenConns(s.cfg.Database.MaxOpenConns)
	// s.db.SetMaxIdleConns(s.cfg.Database.MaxIdleConns)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.Cfg.Database.Timeout)*time.Millisecond)
	defer cancel()

	err = s.Db.Ping(ctx)
	if err != nil {
		panic(fmt.Errorf(`error on verifying db connection %w`, err))
	}
}

func (s *Server) initHandlers(ctx context.Context) {
	// Repositories
	accountRepository := postgres.NewAccountRepo(s.Db)
	transactionRepository := postgres.NewTransactionRepo(s.Db)
	healthRepository := postgres.NewHealthRepo(s.Db)

	// Usecases
	accountUsecase := account.NewUsecase(accountRepository)
	transactionUsecase := transaction.NewUsecase(transactionRepository, accountRepository)
	healthUsecase := health.NewUsecase(healthRepository)

	mux.RegisterAccountHandler(ctx, s.Router, accountUsecase)
	mux.RegisterTransactionHandler(ctx, s.Router, transactionUsecase)
	mux.RegisterErrorHandler(ctx, s.Router)
	mux.RegisterHealthHandler(ctx, s.Router, *healthUsecase)
}

func (s *Server) gracefulShutdown(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	slog.Info("Shutting down...")

	ctx, shutdown := context.WithTimeout(ctx, 3*time.Second)
	defer s.closeResources(ctx, shutdown)

	err := s.HttpServer.Shutdown(ctx)
	if err != nil {
		panic("error on shutting down")
	}
}

func (s *Server) closeResources(ctx context.Context, shutdown context.CancelFunc) {
	s.Db.Conn().Close(ctx)
	shutdown()
}
