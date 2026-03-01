package pubsubapi

import (
	"context"
	"encoding/json"
	"log/slog"

	"cloud.google.com/go/pubsub"
	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewTemplateConsumer)

type TemplateMessage struct {
	ID string `json:"id"`
}

type TemplateConsumer struct {
	client *pubsub.Client
}

func NewTemplateConsumer(client *pubsub.Client) (*TemplateConsumer, error) {
	c := &TemplateConsumer{
		client: client,
	}
	// Note: in a real environment, you probably want to spin this off into a goroutine
	// or manage its lifecycle (so that it doesn't block the IoC initialization).
	// e.g. go c.Start(context.Background())
	return c, nil
}

func (c *TemplateConsumer) Start(ctx context.Context) error {
	sub := c.client.Subscription("template-subscription")

	slog.Info("Starting to listen on subscription", "sub", "template-subscription")

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var payload TemplateMessage
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			slog.Error("error unmarshaling pubsub message", "error", err)
			msg.Nack()
			return
		}

		slog.Info("Successfully processed message", "id", payload.ID)
		// Process business logic here

		msg.Ack()
	})

	return err
}
