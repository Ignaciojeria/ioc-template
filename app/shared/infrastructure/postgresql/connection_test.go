package postgresql

import (
	"strings"
	"testing"

	"archetype/app/shared/configuration"
)

func TestNewConnection_InvalidURL(t *testing.T) {
	conf := configuration.PostgreSQLConfiguration{
		DATABASE_URL:               "invalid_url",
		DATABASE_POSTGRES_HOSTNAME: "localhost", // Triggers the warn override
	}

	db, err := NewConnection(conf)
	if err == nil {
		t.Fatal("expected error connecting with invalid URL, got nil")
	}
	if db != nil {
		t.Errorf("expected nil db on error, got %v", db)
	}

	// This validates that the driver throws an error matching invalid format
	if !strings.Contains(err.Error(), "invalid DSN") {
		t.Errorf("expected DSN formatting error, got %v", err)
	}
}

func TestNewConnection_InvalidDSN(t *testing.T) {
	conf := configuration.PostgreSQLConfiguration{
		DATABASE_URL:               "",
		DATABASE_POSTGRES_USERNAME: "test",
		DATABASE_POSTGRES_PASSWORD: "test",
		DATABASE_POSTGRES_HOSTNAME: "invalid_host",
		DATABASE_POSTGRES_PORT:     "abcd", // invalid port format
		DATABASE_POSTGRES_NAME:     "db",
		DATABASE_POSTGRES_SSL_MODE: "disable",
	}

	db, err := NewConnection(conf)
	if err == nil {
		t.Fatal("expected error connecting with invalid explicit parameters, got nil")
	}
	if db != nil {
		t.Errorf("expected nil db on error, got %v", db)
	}

	if !strings.Contains(err.Error(), "invalid DSN") {
		t.Errorf("expected invalid DSN format error, got %v", err)
	}
}
