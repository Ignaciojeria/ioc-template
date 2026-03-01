package configuration

import (
	"errors"
	"sync"
	"testing"
)

func TestLoadEnvOnce_LoadSuccessBranch(t *testing.T) {
	once = sync.Once{}
	orig := godotenvLoad
	defer func() { godotenvLoad = orig }()

	calls := 0
	godotenvLoad = func(filenames ...string) error {
		calls++
		return nil
	}

	loadEnvOnce()
	loadEnvOnce()

	if calls != 1 {
		t.Fatalf("expected loader to run once, got %d", calls)
	}
}

func TestLoadEnvOnce_LoadErrorBranch(t *testing.T) {
	once = sync.Once{}
	orig := godotenvLoad
	defer func() { godotenvLoad = orig }()

	godotenvLoad = func(filenames ...string) error {
		return errors.New("missing env")
	}

	loadEnvOnce()
}
