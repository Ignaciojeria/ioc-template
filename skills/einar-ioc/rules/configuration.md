# configuration

> Environment configuration with caarlos0/env and godotenv

## app/shared/configuration/conf.go

```go
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
```

---

## app/shared/configuration/parse.go

```go
package configuration

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var once sync.Once

// handleEnvLoad logs the result of loading the .env file. Extracted for testability.
func handleEnvLoad(err error) {
	if err != nil {
		slog.Warn(".env not found, loading environment variables from system.")
	} else {
		slog.Info("Environment variables loaded from .env file.")
	}
}

// loadEnvOnce ensures that the .env file is only loaded once per application lifecycle.
func loadEnvOnce() {
	once.Do(func() {
		handleEnvLoad(godotenv.Load())
	})
}

// Parse loads the .env file (if present) and parses the environment variables into the generic struct T.
// Struct T can use `env:"VAR_NAME"` and `envDefault:"default_value"` tags.
func Parse[T any]() (T, error) {
	loadEnvOnce()
	var conf T
	if err := env.Parse(&conf); err != nil {
		return conf, fmt.Errorf("failed to parse configuration: %w", err)
	}
	return conf, nil
}
```
