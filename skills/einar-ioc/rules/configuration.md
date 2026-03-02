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
	"os"
	"path/filepath"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var once sync.Once

func handleEnvLoad(err error) {
	if err != nil {
		slog.Warn(".env not found, loading environment variables from system.")
	} else {
		slog.Info("Environment variables loaded from .env file.")
	}
}

func findProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	dir := wd
	for dir != filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	return wd
}

func loadEnvOnce() {
	once.Do(func() {
		root := findProjectRoot()
		envPath := filepath.Join(root, ".env")
		handleEnvLoad(godotenv.Load(envPath))
	})
}

func Parse[T any]() (T, error) {
	loadEnvOnce()
	var conf T
	if err := env.Parse(&conf); err != nil {
		return conf, fmt.Errorf("failed to parse configuration: %w", err)
	}
	return conf, nil
}
```

---

## Unit tests

When creating a new component, generate tests following this pattern:

### app/shared/configuration/conf_test.go

```go
package configuration

import (
	"os"
	"strings"
	"testing"

	"archetype"
)

func TestNewConf_DefaultValues(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "conf_test_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	wd, _ := os.Getwd()
	t.Cleanup(func() {
		_ = os.Chdir(wd)
		_ = os.RemoveAll(tmpDir)
		os.Unsetenv("VERSION")
		os.Unsetenv("PORT")
		os.Unsetenv("PROJECT_NAME")
	})

	_ = os.Chdir(tmpDir)
	os.Setenv("VERSION", strings.TrimSpace(archetype.Version))
	os.Unsetenv("PORT")
	os.Unsetenv("PROJECT_NAME")

	conf, err := NewConf()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conf.PORT != "8080" {
		t.Errorf("expected default port 8080, got %s", conf.PORT)
	}
}

func TestNewConf_CustomEnvs(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("PROJECT_NAME", "mytest")
	os.Setenv("VERSION", "2.0")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("PROJECT_NAME")
		os.Unsetenv("VERSION")
	}()

	conf, err := NewConf()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conf.PORT != "9090" {
		t.Errorf("expected port 9090, got %s", conf.PORT)
	}
}
```

---

### app/shared/configuration/parse_test.go

```go
package configuration

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleEnvLoad(t *testing.T) {
	handleEnvLoad(nil)
	handleEnvLoad(errors.New("file not found"))
}

func TestParse(t *testing.T) {
	type Config struct {
		Port int    `env:"TEST_PARSE_PORT" envDefault:"8080"`
		Host string `env:"TEST_PARSE_HOST" envDefault:"localhost"`
	}
	conf, err := Parse[Config]()
	require.NoError(t, err)
	assert.Equal(t, 8080, conf.Port)
	assert.Equal(t, "localhost", conf.Host)
}
```
