package eventbus

import (
	"context"
	"fmt"

	"archetype/app/shared/infrastructure/eventbus"

	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewTemplatePublisher)

// DomainEventPublisher publishes domain events to the broker. Implemented by *TemplatePublisher.
type DomainEventPublisher interface {
	Publish(ctx context.Context, e eventbus.DomainEvent) error
}

type TemplatePublisher struct {
	publisher eventbus.Publisher
}

func NewTemplatePublisher(publisher eventbus.Publisher) (DomainEventPublisher, error) {
	if publisher == nil {
		return nil, fmt.Errorf("publisher dependency is nil")
	}
	return &TemplatePublisher{
		publisher: publisher,
	}, nil
}

func (p *TemplatePublisher) Publish(ctx context.Context, e eventbus.DomainEvent) error {
	request := eventbus.PublishRequest{
		Topic: "your-topic-name",
		Event: e,
	}
	return p.publisher.Publish(ctx, request)
}
