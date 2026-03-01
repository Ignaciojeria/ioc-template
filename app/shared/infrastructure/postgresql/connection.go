package postgresql

import (
	"embed"
	"fmt"
	"log/slog"

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
	// 1️⃣ Si DATABASE_URL no está seteado, armarlo manualmente
	if dsn == "" {
		dsn = "postgres://" + env.DATABASE_POSTGRES_USERNAME + ":" +
			env.DATABASE_POSTGRES_PASSWORD + "@" +
			env.DATABASE_POSTGRES_HOSTNAME + ":" +
			env.DATABASE_POSTGRES_PORT + "/" +
			env.DATABASE_POSTGRES_NAME + "?sslmode=" +
			env.DATABASE_POSTGRES_SSL_MODE
	} else if env.DATABASE_POSTGRES_HOSTNAME != "" {
		// Warning elegante: si hay tanto URL como variables individuales
		slog.Warn("[config warning] DATABASE_URL is set and overrides individual Postgres variables")
	}

	// 2️⃣ Conectar con el driver nativo puro pgx
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// 3️⃣ Correr migraciones automáticamente
	if err := runMigrations(db, env.DATABASE_POSTGRES_NAME); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sqlx.DB, dbName string) error {
	d, err := iofs.New(migrationsFS, "migrations")
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
