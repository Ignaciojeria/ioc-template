package configuration

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleEnvLoad(t *testing.T) {
	t.Run("logs info when no error", func(t *testing.T) {
		handleEnvLoad(nil)
	})

	t.Run("logs warn when error", func(t *testing.T) {
		handleEnvLoad(errors.New("file not found"))
	})
}

func TestParse(t *testing.T) {
	t.Run("parses configuration successfully", func(t *testing.T) {
		type Config struct {
			Port int    `env:"TEST_PARSE_PORT" envDefault:"8080"`
			Host string `env:"TEST_PARSE_HOST" envDefault:"localhost"`
		}

		// Use defaults
		conf, err := Parse[Config]()
		require.NoError(t, err)
		assert.Equal(t, 8080, conf.Port)
		assert.Equal(t, "localhost", conf.Host)

		// Override with env vars
		t.Setenv("TEST_PARSE_PORT", "3000")
		t.Setenv("TEST_PARSE_HOST", "0.0.0.0")
		conf, err = Parse[Config]()
		require.NoError(t, err)
		assert.Equal(t, 3000, conf.Port)
		assert.Equal(t, "0.0.0.0", conf.Host)
	})

	t.Run("returns error when parsing fails", func(t *testing.T) {
		type BadConfig struct {
			Port int `env:"TEST_PARSE_BAD_PORT" envDefault:"8080"`
		}

		t.Setenv("TEST_PARSE_BAD_PORT", "not-a-number")
		_, err := Parse[BadConfig]()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse configuration")
	})

	t.Run("returns error when required field is missing", func(t *testing.T) {
		type RequiredConfig struct {
			APIKey string `env:"TEST_PARSE_REQUIRED_API_KEY,required"`
		}

		os.Unsetenv("TEST_PARSE_REQUIRED_API_KEY")
		_, err := Parse[RequiredConfig]()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse configuration")
	})
}
