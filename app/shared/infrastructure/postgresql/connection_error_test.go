package postgresql

import (
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"testing"

	"archetype/app/shared/configuration"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

type fakeMigrator struct{ err error }

func (f fakeMigrator) Up() error { return f.err }

func TestNewConnection_InvalidURLAfterConnect(t *testing.T) {
	origConnect := sqlxConnect
	origRun := internalRunMigrations
	defer func() {
		sqlxConnect = origConnect
		internalRunMigrations = origRun
	}()

	rawDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected sqlmock error: %v", err)
	}
	db := sqlx.NewDb(rawDB, "sqlmock")
	defer db.Close()

	sqlxConnect = func(driverName, dataSourceName string) (*sqlx.DB, error) { return db, nil }
	internalRunMigrations = func(db *sqlx.DB, dbName string, fsys embed.FS) error { return nil }

	_, err = NewConnection(configuration.Conf{DATABASE_URL: "postgres://user:pass@host/%-invalid"})
	if err == nil {
		t.Fatal("expected invalid url error")
	}
}

func TestNewConnection_MigrationError(t *testing.T) {
	origConnect := sqlxConnect
	origRun := internalRunMigrations
	defer func() {
		sqlxConnect = origConnect
		internalRunMigrations = origRun
	}()

	rawDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected sqlmock error: %v", err)
	}
	db := sqlx.NewDb(rawDB, "sqlmock")
	defer db.Close()

	sqlxConnect = func(driverName, dataSourceName string) (*sqlx.DB, error) { return db, nil }
	internalRunMigrations = func(db *sqlx.DB, dbName string, fsys embed.FS) error { return errors.New("migration failed") }

	_, err = NewConnection(configuration.Conf{DATABASE_URL: "postgres://user:pass@host/db"})
	if err == nil {
		t.Fatal("expected migration error")
	}
}

func TestNewConnection_SuccessWithoutDocker(t *testing.T) {
	origConnect := sqlxConnect
	origRun := internalRunMigrations
	defer func() {
		sqlxConnect = origConnect
		internalRunMigrations = origRun
	}()

	rawDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected sqlmock error: %v", err)
	}
	db := sqlx.NewDb(rawDB, "sqlmock")
	defer db.Close()

	sqlxConnect = func(driverName, dataSourceName string) (*sqlx.DB, error) { return db, nil }
	internalRunMigrations = func(db *sqlx.DB, dbName string, fsys embed.FS) error { return nil }

	got, err := NewConnection(configuration.Conf{DATABASE_URL: "postgres://user:pass@host/db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != db {
		t.Fatal("expected same db from mocked connect")
	}
}

func TestInternalRunMigrations_IOFSError(t *testing.T) {
	origIOFS := iofsNew
	defer func() { iofsNew = origIOFS }()

	iofsNew = func(fsys fs.FS, path string) (source.Driver, error) {
		return nil, errors.New("iofs error")
	}

	rawDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected sqlmock error: %v", err)
	}
	db := sqlx.NewDb(rawDB, "sqlmock")
	defer db.Close()

	err = internalRunMigrations(db, "db", migrationsFS)
	if err == nil {
		t.Fatal("expected iofs error")
	}
}

func TestInternalRunMigrations_MigrateNewError(t *testing.T) {
	origIOFS := iofsNew
	origWithInstance := postgresWithInstance
	origMigrateNew := migrateNewWithInstanceFn
	defer func() {
		iofsNew = origIOFS
		postgresWithInstance = origWithInstance
		migrateNewWithInstanceFn = origMigrateNew
	}()

	rawDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected sqlmock error: %v", err)
	}
	db := sqlx.NewDb(rawDB, "sqlmock")
	defer db.Close()

	iofsNew = func(fsys fs.FS, path string) (source.Driver, error) { return iofs.New(migrationsFS, "migrations") }
	postgresWithInstance = func(instance *sql.DB, config *postgres.Config) (database.Driver, error) { return nil, nil }
	migrateNewWithInstanceFn = func(sourceName string, d source.Driver, databaseName string, driver database.Driver) (migrationUp, error) {
		return nil, errors.New("new migrate error")
	}

	err = internalRunMigrations(db, "db", migrationsFS)
	if err == nil {
		t.Fatal("expected migrate.NewWithInstance error")
	}
}

func TestInternalRunMigrations_UpErrorAndNoChange(t *testing.T) {
	origIOFS := iofsNew
	origWithInstance := postgresWithInstance
	origMigrateNew := migrateNewWithInstanceFn
	defer func() {
		iofsNew = origIOFS
		postgresWithInstance = origWithInstance
		migrateNewWithInstanceFn = origMigrateNew
	}()

	rawDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected sqlmock error: %v", err)
	}
	db := sqlx.NewDb(rawDB, "sqlmock")
	defer db.Close()

	iofsNew = func(fsys fs.FS, path string) (source.Driver, error) { return iofs.New(migrationsFS, "migrations") }
	postgresWithInstance = func(instance *sql.DB, config *postgres.Config) (database.Driver, error) { return nil, nil }

	migrateNewWithInstanceFn = func(sourceName string, d source.Driver, databaseName string, driver database.Driver) (migrationUp, error) {
		return fakeMigrator{err: errors.New("up error")}, nil
	}
	if err = internalRunMigrations(db, "db", migrationsFS); err == nil {
		t.Fatal("expected Up error")
	}

	migrateNewWithInstanceFn = func(sourceName string, d source.Driver, databaseName string, driver database.Driver) (migrationUp, error) {
		return fakeMigrator{err: migrate.ErrNoChange}, nil
	}
	if err = internalRunMigrations(db, "db", migrationsFS); err != nil {
		t.Fatalf("expected ErrNoChange path to succeed, got %v", err)
	}
}
