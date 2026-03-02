package eventbus

import (
	"context"
	"fmt"

	"archetype/app/application/ports/out"
	"archetype/app/shared/infrastructure/eventbus"

	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewTemplatePublisher)

type templatePublisher struct {
	publisher eventbus.Publisher
}

// NewTemplatePublisher returns an implementation of ports/out.DomainEventPublisher.
func NewTemplatePublisher(publisher eventbus.Publisher) (out.DomainEventPublisher, error) {
	if publisher == nil {
		return nil, fmt.Errorf("publisher dependency is nil")
	}
	return &templatePublisher{
		publisher: publisher,
	}, nil
}

func (p *templatePublisher) Publish(ctx context.Context, e out.Event) error {
	domainEvent, ok := e.(eventbus.DomainEvent)
	if !ok {
		return fmt.Errorf("event must implement eventbus.DomainEvent for CloudEvents serialization")
	}
	request := eventbus.PublishRequest{
		Topic: "your-topic-name",
		Event: domainEvent,
	}
	return p.publisher.Publish(ctx, request)
}
