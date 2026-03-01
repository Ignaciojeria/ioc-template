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
	// default to local postgres if not provided by env, excellent for rapid prototyping
	DATABASE_URL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`

	// --- EventBroker Factory ---
	// nats | gcp
	EVENT_BROKER string `env:"EVENT_BROKER" envDefault:"nats"`

	// --- GCP Pub/Sub Configuration (Optional) ---
	GOOGLE_PROJECT_ID string `env:"GOOGLE_PROJECT_ID"`
}

// NewConf loads the configuration and provides it.
// It is returned by value because it's lightweight and immutable.
func NewConf() (Conf, error) {
	return Parse[Conf]()
}
