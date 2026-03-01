---
name: eventbus-messaging
description: Handling GCP and NATS messaging via CloudEvents
---

## Overview

Einar uses a unified `EventBus` architecture based on `cloudevents/sdk-go/v2`. 
Whether the active broker is Google Cloud Pub/Sub or NATS, the implementation is abstracted behind the `eventbus.Publisher` and `eventbus.Subscriber` interfaces.

You MUST always use these interfaces for asynchronous event-driven communication.

## The Interfaces

These interfaces are defined in `app/shared/infrastructure/eventbus`.

> [!NOTE]
> Read `app/adapter/in/eventbus/consumer.go` and `app/adapter/out/eventbus/publisher.go` globally natively within this workspace for living examples of the snippets below.

```go
type Publisher interface {
	Publish(ctx context.Context, request PublishRequest) error
}

type Subscriber interface {
	Start(subscriptionName string, processor MessageProcessor, receiveSettings ReceiveSettings) error
}
```

## Creating a Subscriber

A subscriber lives in `app/adapter/in/eventbus/`. It acts identically to an HTTP controller but for asynchronous events.

### Mandatory Template for Subscribers

```go
package eventbus

import (
	"context"
	"net/http"

	"archetype/app/shared/infrastructure/eventbus"
	"archetype/app/shared/infrastructure/observability"
	"github.com/Ignaciojeria/ioc"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var _ = ioc.Register(NewMyConsumer)

type MyPayload struct {
	ID string `json:"id"`
}

type MyConsumer struct {
	subscriber eventbus.Subscriber
	obs        observability.Observability
	// Inject usecases here
}

func NewMyConsumer(sub eventbus.Subscriber, obs observability.Observability) (*MyConsumer, error) {
	c := &MyConsumer{
		subscriber: sub,
		obs:        obs,
	}
	
	processor := func(ctx context.Context, event cloudevents.Event) int {
		c.obs.Logger.InfoContext(ctx, "CloudEvent received", "id", event.ID(), "type", event.Type())

		var payload MyPayload
		if err := event.DataAs(&payload); err != nil {
			c.obs.Logger.ErrorContext(ctx, "failed_to_unmarshal_cloudevent", "error", err.Error())
			// Invalid payload -> ACK (do not retry infinite loops)
			return http.StatusAccepted
		}

		c.obs.Logger.InfoContext(ctx, "Successfully processed payload", "id", payload.ID)
		// Call core usecase logic here...

		// Returning 200 tells the broker to ACK the message
		return http.StatusOK
	}

	// CRITICAL: Start MUST be called inside a goroutine to prevent blocking the application startup
	go c.subscriber.Start("my-subscription-name", processor, eventbus.ReceiveSettings{MaxOutstandingMessages: 3})
	
	return c, nil
}
```

## Creating a Publisher

A publisher adapter lives in `app/adapter/out/eventbus/`. It implements the outbound port.

```go
package eventbus

import (
	"context"
	"github.com/Ignaciojeria/ioc"
	"archetype/app/shared/infrastructure/eventbus"
)

var _ = ioc.Register(NewMyOutboundPublisher)

type MyOutboundPublisher struct {
	publisher eventbus.Publisher
}

func NewMyOutboundPublisher(p eventbus.Publisher) (*MyOutboundPublisher, error) {
	return &MyOutboundPublisher{publisher: p}, nil
}

// PublishSomething fulfills a `core/usecase` Outbound Port interface
func (p *MyOutboundPublisher) PublishSomething(ctx context.Context, data MyEventData) error {
	return p.publisher.Publish(ctx, eventbus.PublishRequest{
		Topic: "my-topic-name",
		Event: data, // Must implement DomainEvent interface (.ToCloudEvent())
	})
}
```

## Defining Domain Events

Your core domain events must implement `DomainEvent`:

```go
package domain

import cloudevents "github.com/cloudevents/sdk-go/v2"

type MyEventData struct {
	ID string `json:"id"`
}

func (e MyEventData) ToCloudEvent() cloudevents.Event {
	ce := cloudevents.NewEvent()
	ce.SetID(e.ID)
	ce.SetType("my.event.type")
	ce.SetSource("my-service")
	ce.SetData(cloudevents.ApplicationJSON, e)
	return ce
}
```
