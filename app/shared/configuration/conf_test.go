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
