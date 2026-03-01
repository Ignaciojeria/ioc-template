package configuration

import (
	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewConf)

type Conf struct {
	PORT         string `env:"PORT" envDefault:"8080"`
	PROJECT_NAME string `env:"PROJECT_NAME"`
	VERSION      string `env:"VERSION"`

	// --- PostgreSQL Configuration (Optional) ---
	DATABASE_URL               string `env:"DATABASE_URL"`
	DATABASE_POSTGRES_HOSTNAME string `env:"DATABASE_POSTGRES_HOSTNAME"`
	DATABASE_POSTGRES_PORT     string `env:"DATABASE_POSTGRES_PORT"`
	DATABASE_POSTGRES_NAME     string `env:"DATABASE_POSTGRES_NAME"`
	DATABASE_POSTGRES_USERNAME string `env:"DATABASE_POSTGRES_USERNAME"`
	DATABASE_POSTGRES_PASSWORD string `env:"DATABASE_POSTGRES_PASSWORD"`
	DATABASE_POSTGRES_SSL_MODE string `env:"DATABASE_POSTGRES_SSL_MODE"`

	// --- GCP Pub/Sub Configuration (Optional) ---
	GOOGLE_PROJECT_ID string `env:"GOOGLE_PROJECT_ID"`
}

// NewConf loads the configuration and provides it.
// It is returned by value because it's lightweight and immutable.
func NewConf() (Conf, error) {
	return Parse[Conf]()
}
