package eventbus

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"archetype/app/shared/infrastructure/httpserver"

	"cloud.google.com/go/pubsub"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func TestNewGcpSubscriber(t *testing.T) {
	// Dummy Fuego Server (requires it to safely bind)
	// We don't start the server, just pass it for the constructor.
	srv := &httpserver.Server{}

	sub, err := NewGcpSubscriber(nil, srv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub == nil {
		t.Fatal("expected subscriber, got nil")
	}
}

func TestConvertPullMessage(t *testing.T) {
	sub := &GcpSubscriber{}

	// Test case 1: valid CloudEvent in JSON
	ce := cloudevents.NewEvent()
	ce.SetID("101")
	ce.SetType("test.type")
	ce.SetSource("test.source")

	bytesCE, _ := json.Marshal(ce)

	msg := &pubsub.Message{
		ID:   "202", // Message ID should be ignored if CloudEvent has an ID natively
		Data: bytesCE,
	}

	convCE := sub.convertPullMessage("sub-name", msg)
	if convCE.ID() != "101" {
		t.Errorf("expected ID 101, got %s", convCE.ID())
	}

	// Test case 2: invalid JSON (Fallback creation)
	msgInvalid := &pubsub.Message{
		ID:   "303",
		Data: []byte(`{invalid}`),
	}

	fallbackCE := sub.convertPullMessage("sub-name", msgInvalid)
	if fallbackCE.ID() != "303" {
		t.Errorf("expected fallback ID 303, got %s", fallbackCE.ID())
	}
}

func TestMakePushHandler(t *testing.T) {
	sub := &GcpSubscriber{}

	var processedID string
	processor := func(ctx context.Context, event cloudevents.Event) int {
		processedID = event.ID()
		return http.StatusOK
	}

	handler := sub.makePushHandler("sub-name", processor)

	// Test manual custom POST
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer([]byte(`{}`)))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	// Test native GCP Push
	envelopeJSON := `{
		"message": {
			"messageId": "pubsub-id-999",
			"data": "e30=", 
			"attributes": {}
		}
	}`
	reqNative := httptest.NewRequest("POST", "/test", bytes.NewBuffer([]byte(envelopeJSON)))
	reqNative.Header.Set("X-Goog-Channel-ID", "yes")
	wNative := httptest.NewRecorder()

	handler.ServeHTTP(wNative, reqNative)

	if wNative.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", wNative.Code)
	}
	if processedID != "pubsub-id-999" {
		t.Errorf("expected processor to hook event pubsub-id-999, got %s", processedID)
	}
}
