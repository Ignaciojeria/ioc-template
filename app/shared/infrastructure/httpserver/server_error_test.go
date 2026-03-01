package httpserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"archetype/app/shared/configuration"

	"github.com/go-fuego/fuego"
	"github.com/hellofresh/health-go/v5"
)

func TestNewServer_HealthCheckError(t *testing.T) {
	orig := healthNew
	defer func() { healthNew = orig }()

	healthNew = func(opts ...health.Option) (*health.Health, error) {
		return nil, errors.New("boom")
	}

	_, err := NewServer(configuration.Conf{PORT: "8080"})
	if err == nil {
		t.Fatal("expected error when health initialization fails")
	}
}

func TestStartServer_RunError(t *testing.T) {
	origRun := serverRun
	defer func() { serverRun = origRun }()

	serverRun = func(s *fuego.Server) error {
		return errors.New("listen failed")
	}

	s := &Server{Manager: fuego.NewServer(fuego.WithAddr(":0")), conf: configuration.Conf{PORT: "0"}}
	err := StartServer(s)
	if err == nil {
		t.Fatal("expected start error")
	}
}

func TestStartServer_ShutdownErrorPath(t *testing.T) {
	origRun := serverRun
	origStop := serverStop
	origNotify := signalNotify
	defer func() {
		serverRun = origRun
		serverStop = origStop
		signalNotify = origNotify
	}()

	testSignalChan := make(chan os.Signal, 1)
	signalNotify = func(c chan<- os.Signal, sig ...os.Signal) {
		go func() { c <- <-testSignalChan }()
	}

	runReturned := make(chan struct{})
	serverRun = func(s *fuego.Server) error {
		<-runReturned
		return http.ErrServerClosed
	}
	serverStop = func(s *fuego.Server, ctx context.Context) error {
		select {
		case <-runReturned:
		default:
			close(runReturned)
		}
		return errors.New("shutdown failed")
	}

	s := &Server{Manager: fuego.NewServer(fuego.WithAddr(":0")), conf: configuration.Conf{PORT: "0"}}
	errCh := make(chan error, 1)
	go func() { errCh <- StartServer(s) }()

	time.Sleep(50 * time.Millisecond)
	testSignalChan <- syscall.SIGTERM

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting StartServer to return")
	}
}
