package postgresql

import (
	"log/slog"

	"archetype/app/shared/configuration"

	"github.com/Ignaciojeria/ioc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var _ = ioc.Register(NewConnection)

// NewConnection creates a new PostgreSQL GORM connection using the provided configuration.
func NewConnection(env configuration.PostgreSQLConfiguration) (*gorm.DB, error) {

	// 1️⃣ Si DATABASE_URL está seteado -> usarlo sí o sí
	if env.DATABASE_URL != "" {

		// Warning elegante: si también hay definidas variables sueltas
		if env.DATABASE_POSTGRES_HOSTNAME != "" ||
			env.DATABASE_POSTGRES_USERNAME != "" ||
			env.DATABASE_POSTGRES_PASSWORD != "" ||
			env.DATABASE_POSTGRES_NAME != "" {

			slog.Warn("[config warning] DATABASE_URL is set and overrides individual Postgres variables")
		}

		db, err := gorm.Open(postgres.Open(env.DATABASE_URL))
		if err != nil {
			return nil, err
		}

		if err := db.Use(tracing.NewPlugin()); err != nil {
			return nil, err
		}
		return db, nil
	}

	// 2️⃣ Resolver DSN manualmente si no tienes DATABASE_URL
	dsn := "postgres://" + env.DATABASE_POSTGRES_USERNAME + ":" +
		env.DATABASE_POSTGRES_PASSWORD + "@" +
		env.DATABASE_POSTGRES_HOSTNAME + ":" +
		env.DATABASE_POSTGRES_PORT + "/" +
		env.DATABASE_POSTGRES_NAME + "?sslmode=" +
		env.DATABASE_POSTGRES_SSL_MODE

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		return nil, err
	}

	return db, nil
}
