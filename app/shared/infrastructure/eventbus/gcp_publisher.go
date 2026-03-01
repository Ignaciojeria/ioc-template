package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewGcpPublisher)

var jsonMarshal = json.Marshal

type GcpPublisher struct {
	client *pubsub.Client
}

// NewGcpPublisher creates an implementation of the universal Publisher interface backed by GCP Pub/Sub.
func NewGcpPublisher(c *pubsub.Client) (Publisher, error) {
	return &GcpPublisher{client: c}, nil
}

// Publish transforms a DomainEvent into a CloudEvent and dispatches it over GCP pub/sub.
func (p *GcpPublisher) Publish(
	ctx context.Context,
	request PublishRequest,
) error {
	ce := request.Event.ToCloudEvent()

	bytes, err := jsonMarshal(ce)
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
