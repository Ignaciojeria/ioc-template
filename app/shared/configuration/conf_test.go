package configuration

import (
	"os"
	"strings"
	"testing"

	"archetype"
)

func TestNewConf_DefaultValues(t *testing.T) {
	// Let's act like main.go and inject the embedded version
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
	if conf.PROJECT_NAME != "" {
		t.Errorf("expected empty project name, got %s", conf.PROJECT_NAME)
	}
	if conf.VERSION != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", conf.VERSION)
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
	if conf.PROJECT_NAME != "mytest" {
		t.Errorf("expected project name mytest, got %s", conf.PROJECT_NAME)
	}
	if conf.VERSION != "2.0" {
		t.Errorf("expected version 2.0, got %s", conf.VERSION)
	}
}
func TestParse_Error(t *testing.T) {
	os.Setenv("BAD_INT", "not_a_number")
	defer os.Unsetenv("BAD_INT")

	type BadStruct struct {
		Number int `env:"BAD_INT"`
	}
	_, err := Parse[BadStruct]()
	if err == nil {
		t.Error("expected error parsing non-numeric value into int, got nil")
	}
}
