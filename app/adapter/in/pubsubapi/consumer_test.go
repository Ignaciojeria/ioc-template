package pubsubapi

import (
	"testing"
)

func TestNewTemplateConsumer(t *testing.T) {
	// Creating consumer with nil pubsub client just for IoC injection validation.
	// In real tests, interface mocking should be preferred.
	c, err := NewTemplateConsumer(nil)

	if err != nil {
		t.Fatalf("expected no error during consumer creation, got %v", err)
	}

	if c == nil {
		t.Fatal("expected consumer instance, got nil")
	}
}
