package postgresql

import (
	"strings"
	"testing"

	"archetype/app/shared/configuration"
)

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
