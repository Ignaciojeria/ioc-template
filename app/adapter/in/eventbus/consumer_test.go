package eventbus

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"archetype/app/shared/infrastructure/eventbus"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type mockSubscriber struct {
	processor eventbus.MessageProcessor
	mu        sync.Mutex
}

func (m *mockSubscriber) Start(subscriptionName string, processor eventbus.MessageProcessor, receiveSettings eventbus.ReceiveSettings) error {
	m.mu.Lock()
	m.processor = processor
	m.mu.Unlock()
	return nil
}

func (m *mockSubscriber) getProcessor() eventbus.MessageProcessor {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.processor
}

func TestNewTemplateConsumer(t *testing.T) {
	mockSub := &mockSubscriber{}
	c, err := NewTemplateConsumer(mockSub)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if c == nil {
		t.Fatal("expected consumer, got nil")
	}

	// wait max 1 sec for go routine
	for i := 0; i < 100; i++ {
		if mockSub.getProcessor() != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if mockSub.getProcessor() == nil {
		t.Fatal("expected processor to be registered")
	}

	// Test processor with invalid payload
	ce := cloudevents.NewEvent()
	ce.SetID("123")
	ce.SetType("test.type")
	ce.SetSource("test.source")

	// Force invalid JSON via byte array matching expected errors.
	ce.SetData(cloudevents.ApplicationJSON, map[string]interface{}{})
	// Injecting unparsable format
	ce.DataEncoded = []byte(`{"invalid":json}`)

	status := mockSub.getProcessor()(context.Background(), ce)
	if status != http.StatusAccepted {
		t.Errorf("expected status %d for invalid json, got %d", http.StatusAccepted, status)
	}

	// Test processor with valid payload via struct mapping
	ce.SetData(cloudevents.ApplicationJSON, map[string]string{"id": "test-id"})
	status = mockSub.getProcessor()(context.Background(), ce)
	if status != http.StatusOK {
		t.Errorf("expected status %d for valid json, got %d", http.StatusOK, status)
	}
}
