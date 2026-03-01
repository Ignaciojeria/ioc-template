# eventbus-gcp

> GCP Pub/Sub client, publisher, and subscriber implementation

## app/shared/infrastructure/eventbus/gcp_client.go

```go
package eventbus

import (
	"context"
	"errors"
	"fmt"

	"archetype/app/shared/configuration"

	"cloud.google.com/go/pubsub"
	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewGcpClient)

// NewGcpClient creates a new GCP PubSub client using the configuration.
func NewGcpClient(env configuration.Conf) (*pubsub.Client, error) {
	if env.EVENT_BROKER != "gcp" {
		return nil, nil
	}

	if env.GOOGLE_PROJECT_ID == "" {
		return nil, errors.New("GOOGLE_PROJECT_ID is required for PubSub client")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, env.GOOGLE_PROJECT_ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	return client, nil
}
```

---

## app/shared/infrastructure/eventbus/gcp_publisher.go

```go
package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/Ignaciojeria/ioc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var _ = ioc.Register(NewGcpPublisher)

type GcpPublisher struct {
	client *pubsub.Client
}

// NewGcpPublisher creates an implementation of the universal Publisher interface backed by GCP Pub/Sub.
func NewGcpPublisher(c *pubsub.Client) (*GcpPublisher, error) {
	if c == nil {
		return nil, nil
	}
	return &GcpPublisher{client: c}, nil
}

// Publish transforms a DomainEvent into a CloudEvent and dispatches it over GCP pub/sub.
func (p *GcpPublisher) Publish(
	ctx context.Context,
	request PublishRequest,
) error {
	ce := request.Event.ToCloudEvent()

	// Inject OpenTelemetry trace context into the CloudEvent extensions
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	for k, v := range carrier {
		ce.SetExtension(k, v)
	}

	bytes, err := json.Marshal(ce)
	if err != nil {
		return fmt.Errorf("cloudevent marshal error: %w", err)
	}

	// Build attributes for Pub/Sub filtering without needing to deserialize the payload.
	attrs := make(map[string]string)

	if ce.Type() != "" {
		attrs["ce-type"] = ce.Type()
	}
	if ce.Source() != "" {
		attrs["ce-source"] = ce.Source()
	}
	if ce.Subject() != "" {
		attrs["ce-subject"] = ce.Subject()
	}
	if ce.ID() != "" {
		attrs["ce-id"] = ce.ID()
	}

	// Dump CloudEvent extensions so GCP pubsub filtering logic handles it identically.
	for k, v := range ce.Context.GetExtensions() {
		attrs[k] = fmt.Sprintf("%v", v)
	}

	pubTopic := p.client.Topic(request.Topic)
	pubTopic.EnableMessageOrdering = true

	_, err = pubTopic.Publish(ctx, &pubsub.Message{
		Data:        bytes,
		Attributes:  attrs,
		OrderingKey: request.OrderingKey,
	}).Get(ctx)

	return err
}
```

---

## app/shared/infrastructure/eventbus/gcp_subscriber.go

```go
package eventbus

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"archetype/app/shared/infrastructure/httpserver"

	"cloud.google.com/go/pubsub"
	"github.com/Ignaciojeria/ioc"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var _ = ioc.Register(NewGcpSubscriber)

type GcpSubscriber struct {
	client     *pubsub.Client
	httpServer *httpserver.Server
}

func NewGcpSubscriber(c *pubsub.Client, s *httpserver.Server) (*GcpSubscriber, error) {
	if c == nil {
		return nil, nil
	}
	return &GcpSubscriber{client: c, httpServer: s}, nil
}

func (ps *GcpSubscriber) Start(subscriptionName string, processor MessageProcessor, receiveSettings ReceiveSettings) error {
	sub := ps.client.Subscription(subscriptionName)
	sub.ReceiveSettings.MaxOutstandingMessages = receiveSettings.MaxOutstandingMessages

	// PULL consumer running in background
	go func() {
		ctx := context.Background()
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			ce := ps.convertPullMessage(subscriptionName, msg)

			// Extract trace context from CloudEvent extensions
			carrier := propagation.MapCarrier{}
			for k, v := range ce.Extensions() {
				if strVal, ok := v.(string); ok {
					carrier[k] = strVal
				}
			}
			ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

			status := processor(ctx, ce)

			if status >= 500 {
				msg.Nack()
				return
			}
			msg.Ack()
		})

		if err != nil {
			slog.Error("pubsub_receive_failed",
				"subscription", subscriptionName,
				"error", err.Error(),
			)
			time.Sleep(5 * time.Second)
			ps.Start(subscriptionName, processor, receiveSettings) // auto-retry PULL
			return
		}
	}()

	// PUSH consumer via HTTP (Fuego integration)
	path := "/subscription/" + subscriptionName

	fuego.PostStd(ps.httpServer.Manager, path, ps.makePushHandler(subscriptionName, processor), option.Summary("Internal webhook pubsub push"))
	return nil
}

func (ps *GcpSubscriber) convertPullMessage(subName string, msg *pubsub.Message) cloudevents.Event {
	var ce cloudevents.Event
	if err := json.Unmarshal(msg.Data, &ce); err != nil {
		slog.Warn("failed_to_unmarshal_cloudevent",
			"subscription", subName,
			"message_id", msg.ID,
			"error", err.Error(),
		)
		ce = cloudevents.NewEvent()
		ce.SetID(msg.ID)
		ce.SetType("google.pubsub.pull.fallback")
		ce.SetSource("gcp.pubsub/" + subName)
		ce.SetData(cloudevents.ApplicationJSON, msg.Data)
	} else if ce.ID() == "" {
		ce.SetID(msg.ID)
	}
	return ce
}

func (ps *GcpSubscriber) makePushHandler(subName string, processor MessageProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Native GCP push header detection
		if r.Header.Get("X-Goog-Channel-ID") != "" || r.Header.Get("ce-id") != "" {
			ps.handleNativePush(subName, processor, w, r)
			return
		}

		// Manual custom POST for local testing without Cloud Emulator
		ps.handleManualPush(subName, processor, w, r)
	}
}

func (ps *GcpSubscriber) handleNativePush(subName string, processor MessageProcessor, w http.ResponseWriter, r *http.Request) {
	var envelope struct {
		Message struct {
			MessageID  string            `json:"messageId"`
			Data       []byte            `json:"data"`
			Attributes map[string]string `json:"attributes"`
		} `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ce cloudevents.Event
	if err := json.Unmarshal(envelope.Message.Data, &ce); err != nil {
		slog.Warn("failed_to_unmarshal_cloudevent_push",
			"subscription", subName,
			"message_id", envelope.Message.MessageID,
			"error", err.Error(),
		)
		ce = cloudevents.NewEvent()
		ce.SetID(envelope.Message.MessageID)
		ce.SetType("google.pubsub.push.fallback")
		ce.SetSource("gcp.pubsub/" + subName)
		ce.SetData(cloudevents.ApplicationJSON, envelope.Message.Data)
	} else if ce.ID() == "" {
		ce.SetID(envelope.Message.MessageID)
	}

	// Extract trace context
	carrier := propagation.MapCarrier{}
	for k, v := range ce.Extensions() {
		if strVal, ok := v.(string); ok {
			carrier[k] = strVal
		}
	}
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)

	w.WriteHeader(processor(ctx, ce))
}

func (ps *GcpSubscriber) handleManualPush(subName string, processor MessageProcessor, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ce := cloudevents.NewEvent()
	ce.SetID("")
	ce.SetType("manual.message")
	ce.SetSource("manual/" + subName)
	ce.SetData(cloudevents.ApplicationJSON, body)

	for key, values := range r.Header {
		if len(values) > 0 {
			ce.SetExtension(strings.ToLower(key), strings.Join(values, ","))
		}
	}

	// Extract trace context
	carrier := propagation.MapCarrier{}
	for k, v := range ce.Extensions() {
		if strVal, ok := v.(string); ok {
			carrier[k] = strVal
		}
	}
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)

	w.WriteHeader(processor(ctx, ce))
}
```
