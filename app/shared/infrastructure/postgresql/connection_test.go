package postgresql

import (
	"context"
	"strings"
	"testing"
	"time"

	"archetype/app/shared/configuration"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewConnection_Success(t *testing.T) {
	ctx := context.Background()

	// Spin up a PostgreSQL test container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	conf := configuration.Conf{
		DATABASE_URL: connStr,
	}

	// For test purpose, point to dummy but valid FS migrations.
	// Since migrations embed requires files to exist inside package dir,
	// and they exist at `migrations/000001_initial_schema.up.sql`, this will auto-run them natively!

	db, err := NewConnection(conf)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Validate DB ping
	err = db.Ping()
	assert.NoError(t, err)
	db.Close()
}

func TestNewConnection_InvalidURL(t *testing.T) {
	conf := configuration.Conf{
		DATABASE_URL: "invalid_url",
	}

	db, err := NewConnection(conf)
	if err == nil {
		t.Fatal("expected error connecting with invalid URL, got nil")
	}
	if db != nil {
		t.Errorf("expected nil db on error, got %v", db)
	}

	// Validate generic DSN parsing failure
	if !strings.Contains(err.Error(), "failed to connect") {
		t.Errorf("expected connection formatting error, got %v", err)
	}
}

func TestNewConnection_MalformedURL(t *testing.T) {
	// A URL that fails url.Parse
	conf := configuration.Conf{
		// Starting with a colon but no scheme often confuses parser
		DATABASE_URL: "postgres://user:pass@host:port/%-invalid",
	}

	db, err := NewConnection(conf)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if db != nil {
		t.Errorf("expected nil db, got %v", db)
	}
}

func TestNewConnection_EmptyURL(t *testing.T) {
	conf := configuration.Conf{
		DATABASE_URL: "",
	}

	db, err := NewConnection(conf)
	if err == nil {
		t.Fatal("expected error connecting with empty URL, got nil")
	}
	if db != nil {
		t.Errorf("expected nil db on error, got %v", db)
	}

	if !strings.Contains(err.Error(), "DATABASE_URL is not set") {
		t.Errorf("expected connection error, got %v", err)
	}
}
