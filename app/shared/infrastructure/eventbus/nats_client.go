package eventbus

import (
	"log"
	"time"

	"archetype/app/shared/configuration"

	"github.com/Ignaciojeria/ioc"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

var _ = ioc.Register(NewNatsClient)

// NatsClient encapsulates an in-memory NATS server and the connection to it.
// To use a remote server instead, ignore the EmbeddedServer initialization
// and configure nats.Connect with the real URL.
type NatsClient struct {
	EmbeddedServer *server.Server
	Connection     *nats.Conn
}

// NewNatsClient sets up a lightweight in-memory NATS server natively
// within the Go application and connects a client to it. Very useful for local
// development and tests without requiring heavy infrastructure.
func NewNatsClient(conf configuration.Conf) (*NatsClient, error) {

	// 1) Spin up embedded NATS on a random local port or predefined one
	opts := &server.Options{
		// -1 dynamically picks a port, but if we need a static one for
		// explicit monitoring we can define it. -1 is safest for local testing.
		Port: -1,
	}

	embedded, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}

	// 2) Start in background
	go embedded.Start()

	if !embedded.ReadyForConnections(5 * time.Second) {
		return nil, err
	}

	log.Printf("Embedded NATS Server running on %s", embedded.ClientURL())

	// 3) Connect a client locally
	conn, err := nats.Connect(embedded.ClientURL())
	if err != nil {
		return nil, err
	}

	return &NatsClient{
		EmbeddedServer: embedded,
		Connection:     conn,
	}, nil
}
