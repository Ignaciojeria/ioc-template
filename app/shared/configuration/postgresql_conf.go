package configuration

import (
	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewPostgreSQLConfiguration)

type PostgreSQLConfiguration struct {
	DATABASE_URL               string `env:"DATABASE_URL"`
	DATABASE_POSTGRES_HOSTNAME string `env:"DATABASE_POSTGRES_HOSTNAME"`
	DATABASE_POSTGRES_PORT     string `env:"DATABASE_POSTGRES_PORT"`
	DATABASE_POSTGRES_NAME     string `env:"DATABASE_POSTGRES_NAME"`
	DATABASE_POSTGRES_USERNAME string `env:"DATABASE_POSTGRES_USERNAME"`
	DATABASE_POSTGRES_PASSWORD string `env:"DATABASE_POSTGRES_PASSWORD"`
	DATABASE_POSTGRES_SSL_MODE string `env:"DATABASE_POSTGRES_SSL_MODE"`
}

func NewPostgreSQLConfiguration() (PostgreSQLConfiguration, error) {
	return Parse[PostgreSQLConfiguration]()
}
