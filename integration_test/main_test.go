package test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/adrianuf22/back-test-psmo/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var pgContainer *postgres.PostgresContainer

func TestMain(m *testing.M) {
	config.New()
	ctx := context.Background()

	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	var err error
	pgContainer, err = postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:alpine"),
		postgres.WithInitScripts(
			filepath.Join("..", "misc", "database", "0-init.sql"),
			filepath.Join(".", "fixtures", "account_fixtures.sql"),
		),
		postgres.WithDatabase("test_psmo"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test1234"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	s, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	os.Setenv("POSTGRES_DSN", s)

	code := m.Run()

	if err := pgContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate pgContainer: %s", err)
	}

	os.Exit(code)
}
