package httpserver

import (
	"context"
	"testing"
	"time"

	"archetype/app/shared/configuration"
)

func TestNewServer(t *testing.T) {
	conf := configuration.Conf{
		PORT:         "8081",
		PROJECT_NAME: "test-project",
		VERSION:      "v1",
	}

	server, err := NewServer(conf)
	if err != nil {
		t.Fatalf("unexpected error creating server: %v", err)
	}

	if server.Manager == nil {
		t.Fatal("expected server.Manager to be initialized")
	}
}

func TestStartServer(t *testing.T) {
	conf := configuration.Conf{
		PORT:         "0",
		PROJECT_NAME: "test-start",
		VERSION:      "v1",
	}

	server, err := NewServer(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	errChan := make(chan error, 1)

	// StartServer is blocking, so we run it in a goroutine
	go func() {
		errChan <- StartServer(server)
	}()

	// Give the server a brief moment to start up
	time.Sleep(200 * time.Millisecond)

	// Since we are running in tests, we can manually trigger a shutdown
	// instead of sending OS signals (which can be flaky on Windows).
	// Calling Shutdown forces the internal .Run() to exit cleanly.
	if err := server.Manager.Shutdown(context.Background()); err != nil {
		t.Fatalf("failed to shutdown test server: %v", err)
	}

	// Wait for StartServer to return
	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("StartServer returned an unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("StartServer took too long to return after shutdown")
	}
}

func TestStartServer_InvalidPort(t *testing.T) {
	conf := configuration.Conf{
		PORT:         "-1",
		PROJECT_NAME: "test-port",
		VERSION:      "v1",
	}

	server, err := NewServer(conf)
	if err != nil {
		t.Fatalf("unexpected error creating server: %v", err)
	}

	// This should return an error synchronously because the server cannot start on port -1
	err = StartServer(server)
	if err == nil {
		t.Fatal("expected error due to invalid port, got nil")
	}
}
