# postgresql-connection

> PostgreSQL connection with embedded migrations via sqlx and golang-migrate

## app/shared/infrastructure/postgresql/connection.go

```go
package postgresql

import (
	"embed"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"archetype/app/shared/configuration"

	"github.com/Ignaciojeria/ioc"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // register pgx driver
	"github.com/jmoiron/sqlx"
)

var _ = ioc.Register(NewConnection)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// NewConnection creates a new PostgreSQL sqlx connection using the provided configuration.
// It automatically executes any pending migrations encoded in the migrationsFS embedded folder.
func NewConnection(env configuration.Conf) (*sqlx.DB, error) {

	dsn := env.DATABASE_URL
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	// 1️⃣ Conectar con el driver nativo puro pgx
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// 2️⃣ Extraer nombre de la base de datos para las migraciones
	u, err := url.Parse(dsn)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("invalid DATABASE_URL format: %w", err)
	}
	dbName := strings.TrimPrefix(u.Path, "/")

	// 3️⃣ Correr migraciones automáticamente
	if err := internalRunMigrations(db, dbName, migrationsFS); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func internalRunMigrations(db *sqlx.DB, dbName string, fsys embed.FS) error {
	if db == nil {
		return fmt.Errorf("db connection is nil")
	}
	d, err := iofs.New(fsys, "migrations")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		d,
		dbName,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	slog.Info("Database migrations validated/applied successfully")
	return nil
}

// Deprecated: use NewConnection which handles migrations internally.
// Function signature kept for backward compatibility if needed by generated code.
func runMigrations(db *sqlx.DB, dbName string) error {
	return internalRunMigrations(db, dbName, migrationsFS)
}
```
