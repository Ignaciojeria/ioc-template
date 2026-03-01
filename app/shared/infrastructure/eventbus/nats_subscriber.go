package eventbus

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"archetype/app/shared/infrastructure/httpserver"

	"github.com/Ignaciojeria/ioc"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var _ = ioc.Register(NewNatsSubscriber)

// NatsSubscriber implements the Subscriber contract fetching from NATS locally.
type NatsSubscriber struct {
	client        *NatsClient
	server        *httpserver.Server
	subscriptions []*nats.Subscription
}

// NewNatsSubscriber initializes the struct using the in-memory NATS client
func NewNatsSubscriber(client *NatsClient, srv *httpserver.Server) (*NatsSubscriber, error) {
	if client == nil {
		return nil, nil // Disabled via configuration
	}
	return &NatsSubscriber{
		client: client,
		server: srv,
	}, nil
}

// Start adds a subscription natively inside the NATS broker
func (s *NatsSubscriber) Start(subscriptionName string, processor MessageProcessor, receiveSettings ReceiveSettings) error {
	// Create a NATS local QueueSubscription for load-balancing semantics if needed, or normal sub.
	sub, err := s.client.Connection.Subscribe(subscriptionName, func(m *nats.Msg) {
		s.processMessageAsCloudEvent(m, processor)
	})

	if err != nil {
		return err
	}

	s.subscriptions = append(s.subscriptions, sub)
	return nil
}

// processMessageAsCloudEvent unwraps the NATS payload exactly like GCP PULL wrappers do.
func (s *NatsSubscriber) processMessageAsCloudEvent(m *nats.Msg, processor MessageProcessor) {
	ctx := context.Background()

	var ce cloudevents.Event
	err := json.Unmarshal(m.Data, &ce)

	if err != nil {
		// Log errors similar to how we handled Google Cloud Events
		log.Printf("invalid cloudevent on NATS: %v", err)
		return
	}

	// Unpack OTel propagating context
	carrier := propagation.MapCarrier{}
	for k, v := range ce.Extensions() {
		if strVal, ok := v.(string); ok {
			carrier[k] = strVal
		}
	}
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	// In memory, we act like handlers responding 200 strictly.
	statusCode := processor(ctx, ce)
	if statusCode != http.StatusOK {
		log.Printf("Processor failed NATS message with status %d", statusCode)
	}
}
