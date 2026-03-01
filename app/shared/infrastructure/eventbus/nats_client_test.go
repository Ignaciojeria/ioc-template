package eventbus

import (
	"context"
	"net/http"
	"testing"
	"time"

	"archetype/app/shared/configuration"
	"archetype/app/shared/infrastructure/httpserver"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

// NatsDummyEvent for test integration
type NatsDummyEvent struct {
	Message string `json:"message"`
}

func (d NatsDummyEvent) ToCloudEvent() cloudevents.Event {
	e := cloudevents.NewEvent()
	e.SetID("test-id")
	e.SetType("test.dummy.event")
	e.SetData(cloudevents.ApplicationJSON, d)
	return e
}

func TestNatsIntegrationSuite(t *testing.T) {
	conf := configuration.Conf{
		PORT:         "0",
		PROJECT_NAME: "test",
		VERSION:      "1.0",
	}

	// 1) Initialize the Embedded Cluster
	client, err := NewNatsClient(conf)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	defer client.EmbeddedServer.Shutdown()
	defer client.Connection.Close()

	// 2) Prepare the tools
	pub := NewNatsPublisher(client)
	srv := &httpserver.Server{}
	sub, _ := NewNatsSubscriber(client, srv)

	// 3) Hook up a handler capturing success
	received := make(chan bool, 1)

	sub.Register("test-topic", func(ctx context.Context, e cloudevents.Event) int {
		var payload NatsDummyEvent
		if err := e.DataAs(&payload); err == nil {
			if payload.Message == "hello from memory" {
				received <- true
				return http.StatusOK
			}
		}
		return http.StatusBadRequest
	})

	// 4) Publish to Memory
	event := NatsDummyEvent{Message: "hello from memory"}
	req := PublishRequest{
		Topic: "test-topic",
		Event: event,
	}

	err = pub.Publish(context.Background(), req)
	assert.NoError(t, err)

	// 5) Verify asynchronous delivery natively inside Go
	select {
	case <-received:
		// Success!
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for NATS message")
	}
}
