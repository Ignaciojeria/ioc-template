package eventbus

import (
	"context"
	"encoding/json"

	"github.com/Ignaciojeria/ioc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var _ = ioc.Register(NewNatsPublisher)

// NatsPublisher implements Publisher matching GCP semantics but over NATS core
type NatsPublisher struct {
	client *NatsClient
}

// NewNatsPublisher creates a Publisher using the Embedded NATS client connection.
func NewNatsPublisher(client *NatsClient) *NatsPublisher {
	if client == nil {
		return nil
	}
	return &NatsPublisher{
		client: client,
	}
}

// Publish takes a DomainEvent, converts it to a standard CloudEvent, serializes it,
// grabs the OTel active span, injects it into CloudEvent Extensions, and
// finally publishes it over the embedded NATS broker.
func (p *NatsPublisher) Publish(ctx context.Context, request PublishRequest) error {

	// 1) Translate domain to CloudEvent wrapper
	ce := request.Event.ToCloudEvent()
	if ce.ID() == "" {
		ce.SetID("nats-" + request.Topic)
	}
	ce.SetSource("ioc-service")

	// 2) Inject OpenTelemetry context for distributed tracing continuity
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	for k, v := range carrier {
		ce.SetExtension(k, v)
	}

	// 3) Serialize CloudEvent object strictly to JSON bytes
	body, err := json.Marshal(ce)
	if err != nil {
		return err
	}

	// 4) Publish down into NATS subject directly
	return p.client.Connection.Publish(request.Topic, body)
}
