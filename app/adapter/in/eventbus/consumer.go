package eventbus

import (
	"context"
	"log/slog"
	"net/http"

	"archetype/app/shared/infrastructure/eventbus"

	"github.com/Ignaciojeria/ioc"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var _ = ioc.Register(NewTemplateConsumer)

type TemplateMessage struct {
	ID string `json:"id"`
}

type TemplateConsumer struct {
	subscriber eventbus.Subscriber
}

func NewTemplateConsumer(sub eventbus.Subscriber) (*TemplateConsumer, error) {
	c := &TemplateConsumer{
		subscriber: sub,
	}

	processor := func(ctx context.Context, event cloudevents.Event) int {
		slog.Info("CloudEvent received", "id", event.ID(), "type", event.Type())

		var payload TemplateMessage
		if err := event.DataAs(&payload); err != nil {
			slog.Error("failed_to_unmarshal_cloudevent", "error", err.Error())
			// Invalid payload -> ACK (do not retry infinite loops)
			return http.StatusAccepted
		}

		slog.Info("Successfully processed payload", "id", payload.ID)
		// Process core business logic here...

		// Returning 200 tells the broker to ACK the message
		return http.StatusOK
	}

	// This starts the listening process in background via PULL and binds the PUSH http route.
	// You might want to parameterize "template_topic_or_hook" with an environment variable.
	go c.subscriber.Start("template_topic_or_hook", processor, eventbus.ReceiveSettings{MaxOutstandingMessages: 3})

	return c, nil
}
