package configuration

import (
	"os"
	"testing"
)

func TestNewPostgreSQLConfiguration(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")
	os.Setenv("DATABASE_POSTGRES_HOSTNAME", "localhost")

	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("DATABASE_POSTGRES_HOSTNAME")
	}()

	conf, err := NewPostgreSQLConfiguration()
	if err != nil {
		t.Fatalf("unexpected error parsing postgres config: %v", err)
	}

	if conf.DATABASE_URL != "postgres://test:test@localhost:5432/test?sslmode=disable" {
		t.Errorf("expected overriding DATABASE_URL got %s", conf.DATABASE_URL)
	}
	if conf.DATABASE_POSTGRES_HOSTNAME != "localhost" {
		t.Errorf("expected localhost got %s", conf.DATABASE_POSTGRES_HOSTNAME)
	}
}
