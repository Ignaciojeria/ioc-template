package out

import "context"

// Event is a domain event to be published.
// Implementations convert it to the broker format (e.g. CloudEvents) in the adapter.
type Event interface {
	EventType() string
}

// DomainEventPublisher publishes domain events to a broker.
// Implementations (NATS, GCP Pub/Sub, Kafka) live in adapter/out.
type DomainEventPublisher interface {
	Publish(ctx context.Context, e Event) error
}
