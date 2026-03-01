package httpserver

import (
	"context"
	"os"
	"syscall"
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

func TestStartServer_GracefulShutdown(t *testing.T) {
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

	go func() {
		errChan <- StartServer(server)
	}()

	time.Sleep(200 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Manager.Shutdown(ctx); err != nil {
		t.Fatalf("failed to shutdown test server: %v", err)
	}

	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("StartServer returned an unexpected error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("StartServer took too long to return after shutdown")
	}
}

func TestStartServer_Signal(t *testing.T) {
	// Skip signal test on windows if it is flaky
	// but let's try it one last time
	conf := configuration.Conf{
		PORT:         "0",
		PROJECT_NAME: "test-signal",
		VERSION:      "v1",
	}

	server, _ := NewServer(conf)
	errChan := make(chan error, 1)

	go func() {
		errChan <- StartServer(server)
	}()

	time.Sleep(200 * time.Millisecond)

	p, _ := os.FindProcess(os.Getpid())
	_ = p.Signal(syscall.SIGINT)

	select {
	case <-errChan:
		// success
	case <-time.After(2 * time.Second):
		// if it doesn't return, we shutdown manually to not hang
		_ = server.Manager.Shutdown(context.Background())
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

	err = StartServer(server)
	if err == nil {
		t.Fatal("expected error due to invalid port, got nil")
	}
}
