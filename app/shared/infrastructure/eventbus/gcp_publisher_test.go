package eventbus

import (
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
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

func TestNewGcpPublisher(t *testing.T) {
	pub, err := NewGcpPublisher(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pub == nil {
		t.Fatal("expected publisher, got nil")
	}
}

// Full publish test requires a real or pstest client, so we skip execution
// or test just the payload marshalling if we refactored it.
// For now, we only cover the constructor.
