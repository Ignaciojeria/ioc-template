package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/hellofresh/health-go/v5"

	"archetype/app/shared/configuration"
)

var (
	_ = ioc.Register(NewServer)
	_ = ioc.RegisterAtEnd(StartServer)

	healthNew    = health.New
	serverRun    = func(s *fuego.Server) error { return s.Run() }
	serverStop   = func(s *fuego.Server, ctx context.Context) error { return s.Shutdown(ctx) }
	signalNotify = signal.Notify
)

type Server struct {
	Manager *fuego.Server
	conf    configuration.Conf
}

// NewServer creates a new instance of the HTTP Fuego Server.
// It returns a pointer because it manages network state.
func NewServer(conf configuration.Conf) (*Server, error) {
	s := fuego.NewServer(fuego.WithAddr(":" + conf.PORT))

	server := &Server{
		Manager: s,
		conf:    conf,
	}

	if err := server.healthCheck(); err != nil {
		return nil, fmt.Errorf("failed to init healthcheck: %w", err)
	}

	return server, nil
}

// StartServer runs at the end of the dependency graph and starts the HTTP server.
// It blocks the main thread and gracefully handles OS signals.
func StartServer(s *Server) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signalNotify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Second*5)
		defer shutdownCancel()

		fmt.Println("Shutting down server gracefully...")
		if err := serverStop(s.Manager, shutdownCtx); err != nil {
			fmt.Printf("Shutdown error: %v\n", err)
		}
		cancel()
	}()

	fmt.Printf("Starting HTTP server on port %s\n", s.conf.PORT)
	if err := serverRun(s.Manager); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

func (s *Server) healthCheck() error {
	h, err := healthNew(
		health.WithComponent(health.Component{
			Name:    s.conf.PROJECT_NAME,
			Version: s.conf.VERSION,
		}),
		health.WithSystemInfo(),
	)
	if err != nil {
		return err
	}

	fuego.GetStd(s.Manager,
		"/health",
		h.Handler().ServeHTTP,
		option.Summary("healthCheck"),
	)
	return nil
}
