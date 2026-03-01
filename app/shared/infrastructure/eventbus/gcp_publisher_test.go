package eventbus

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DummyEvent struct{}

func (d DummyEvent) ToCloudEvent() cloudevents.Event {
	ce := cloudevents.NewEvent()
	ce.SetID("123")
	ce.SetType("test.type")
	ce.SetSource("test.source")
	ce.SetData(cloudevents.ApplicationJSON, map[string]string{"foo": "bar"})
	ce.SetExtension("customext", "value")
	return ce
}

type MinimalEvent struct{}

func (m MinimalEvent) ToCloudEvent() cloudevents.Event {
	ce := cloudevents.NewEvent()
	// No ID, No Source, No Type explicitly set by user (SDK might set some defaults)
	return ce
}

type FullEvent struct{}

func (f FullEvent) ToCloudEvent() cloudevents.Event {
	ce := cloudevents.NewEvent()
	ce.SetID("full-id")
	ce.SetType("full.type")
	ce.SetSource("full.source")
	ce.SetSubject("full.subject")
	return ce
}

func TestNewGcpPublisher(t *testing.T) {
	pub, err := NewGcpPublisher(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pub == nil {
		t.Fatal("expected publisher, got nil")
	}
}

func TestGcpPublisher_Publish(t *testing.T) {
	ctx := context.Background()

	// Start a fake PubSub server
	srv := pstest.NewServer()
	defer srv.Close()

	// Connect to it securely via gRPC
	conn, err := grpc.NewClient(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial test server: %v", err)
	}
	defer conn.Close()

	// Create test client attached to the fake server
	client, err := pubsub.NewClient(ctx, "project-id", option.WithGRPCConn(conn))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Create a topic on the fake server
	_, err = client.CreateTopic(ctx, "test-topic")
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}

	pub, _ := NewGcpPublisher(client)

	req := PublishRequest{
		Topic:       "test-topic",
		OrderingKey: "123-group",
		Event:       DummyEvent{},
	}

	err = pub.Publish(ctx, req)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// Wait and check if message was sent to the server correctly
	msgs := srv.Messages()
	if len(msgs) == 0 {
		t.Fatalf("No messages published to fake server")
	}
	msg := msgs[0]

	if msg.Attributes["ce-id"] != "123" {
		t.Errorf("Expected ce-id=123, got %s", msg.Attributes["ce-id"])
	}

	if msg.Attributes["ce-type"] != "test.type" {
		t.Errorf("Expected ce-type=test.type, got %s", msg.Attributes["ce-type"])
	}

	if msg.Attributes["customext"] != "value" {
		t.Errorf("Expected customext=value, got %s", msg.Attributes["customext"])
	}

	// Test Minimal Event
	_ = pub.Publish(ctx, PublishRequest{Topic: "test-topic", Event: MinimalEvent{}})

	// Test Full Event
	_ = pub.Publish(ctx, PublishRequest{Topic: "test-topic", Event: FullEvent{}})

	// Test Publish error - cancelled context causes Get() to fail
	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel()
	err = pub.Publish(cancelledCtx, PublishRequest{Topic: "test-topic", Event: DummyEvent{}})
	if err == nil {
		t.Error("expected error when context is cancelled")
	}

	// Check Full Event was received - find by ce-id
	msgs = srv.Messages()
	for _, m := range msgs {
		if m.Attributes["ce-id"] == "full-id" {
			if m.Attributes["ce-subject"] != "full.subject" {
				t.Errorf("Expected ce-subject=full.subject, got %s", m.Attributes["ce-subject"])
			}
			break
		}
	}
}
