package configuration

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var once sync.Once

// loadEnvOnce ensures that the .env file is only loaded once per application lifecycle.
func loadEnvOnce() {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			slog.Warn(".env not found, loading environment variables from system.")
		} else {
			slog.Info("Environment variables loaded from .env file.")
		}
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
