package eventbus

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"archetype/app/shared/configuration"
	"archetype/app/shared/infrastructure/httpserver"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNewGcpSubscriber(t *testing.T) {
	srv := &httpserver.Server{}

	sub, err := NewGcpSubscriber(nil, srv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub == nil {
		t.Fatal("expected subscriber, got nil")
	}
}

func TestGcpSubscriber_Start(t *testing.T) {
	ctx := context.Background()

	conf := configuration.Conf{
		PORT:         "0",
		PROJECT_NAME: "test",
		VERSION:      "v1",
	}
	srv, err := httpserver.NewServer(conf)
	if err != nil {
		t.Fatalf("failed to create fake http server: %v", err)
	}

	testSrv := pstest.NewServer()
	defer testSrv.Close()

	conn, err := grpc.NewClient(testSrv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial test server: %v", err)
	}
	defer conn.Close()

	client, err := pubsub.NewClient(ctx, "project-id", option.WithGRPCConn(conn))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	topic, err := client.CreateTopic(ctx, "test-topic")
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}
	_, err = client.CreateSubscription(ctx, "test-sub", pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		t.Fatalf("failed to create sub: %v", err)
	}

	subscriber, _ := NewGcpSubscriber(client, srv)

	receivedChan := make(chan bool)

	processor := func(ctx context.Context, event cloudevents.Event) int {
		receivedChan <- true
		if event.ID() == "fail-me" {
			return http.StatusInternalServerError // trigger a Nack
		}
		return http.StatusOK // trigger an Ack
	}

	err = subscriber.Start("test-sub", processor, ReceiveSettings{MaxOutstandingMessages: 1})
	if err != nil {
		t.Fatalf("Subscriber start failed: %v", err)
	}

	// Publish success event
	ce1 := cloudevents.NewEvent()
	ce1.SetID("1")
	ce1.SetType("test")
	ce1.SetSource("test")
	data1, _ := json.Marshal(ce1)

	topic.Publish(ctx, &pubsub.Message{Data: data1})
	select {
	case <-receivedChan:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for pubsub pull message")
	}

	// Publish fail event
	ce2 := cloudevents.NewEvent()
	ce2.SetID("fail-me")
	ce2.SetType("test")
	ce2.SetSource("test")
	data2, _ := json.Marshal(ce2)

	topic.Publish(ctx, &pubsub.Message{Data: data2})
	select {
	case <-receivedChan:
		// success reaching the processor (nack will just return to pubsub internally)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for pubsub fail message")
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

	// Test case 3: JSON without ID
	ceNoID := cloudevents.NewEvent()
	ceNoID.SetType("test")
	ceNoID.SetSource("test")
	bytesNoID, _ := json.Marshal(ceNoID)

	msgNoID := &pubsub.Message{
		ID:   "505",
		Data: bytesNoID,
	}

	fallbackNoID := sub.convertPullMessage("sub-name", msgNoID)
	if fallbackNoID.ID() != "505" {
		t.Errorf("expected fallback ID 505 for missing id, got %s", fallbackNoID.ID())
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

	// Test manual custom POST with Headers to check if mapping works
	reqHeaders := httptest.NewRequest("POST", "/test", bytes.NewBuffer([]byte(`{}`)))
	reqHeaders.Header.Set("Custom-Ext", "works")
	wHeaders := httptest.NewRecorder()

	handler.ServeHTTP(wHeaders, reqHeaders)

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

	// Test native GCP Push Invalid
	invalidNativeReq := httptest.NewRequest("POST", "/test", bytes.NewBuffer([]byte(`{invalid}`)))
	invalidNativeReq.Header.Set("X-Goog-Channel-ID", "yes")
	wInvalid := httptest.NewRecorder()

	handler.ServeHTTP(wInvalid, invalidNativeReq)
	if wInvalid.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for invalid json block, got %d", wInvalid.Code)
	}

	// Test native GCP Push fallback data json
	fallbackDataJSON := `{
		"message": {
			"messageId": "pubsub-id-fallback",
			"data": "aW52YWxpZCBqc29uIG9uIHJlc3BvbnNl", 
			"attributes": {}
		}
	}`
	fallbackReq := httptest.NewRequest("POST", "/test", bytes.NewBuffer([]byte(fallbackDataJSON)))
	fallbackReq.Header.Set("X-Goog-Channel-ID", "yes")
	wFallback := httptest.NewRecorder()

	handler.ServeHTTP(wFallback, fallbackReq)
	if processedID != "pubsub-id-fallback" {
		t.Errorf("expected processor to fallback ID, got %s", processedID)
	}

	// Test native GCP Push fallback ID (no native ID in data)
	noIDCE := cloudevents.NewEvent()
	noIDCE.SetType("test")
	noIDCE.SetSource("http://test.source")
	bytesNoID, _ := json.Marshal(noIDCE)
	// Base64 encode it for envelope
	encodedNoID := base64.StdEncoding.EncodeToString(bytesNoID)

	noIDEnvelope := fmt.Sprintf(`{
		"message": {
			"messageId": "mapped-id",
			"data": "%s", 
			"attributes": {}
		}
	}`, encodedNoID)

	reqNoID := httptest.NewRequest("POST", "/test", bytes.NewBuffer([]byte(noIDEnvelope)))
	reqNoID.Header.Set("X-Goog-Channel-ID", "yes")
	wNoID := httptest.NewRecorder()

	handler.ServeHTTP(wNoID, reqNoID)
	if processedID != "mapped-id" {
		t.Errorf("expected ID mapped-id, got %s", processedID)
	}

	// Test Push reading manual body error
	reqErrReader := httptest.NewRequest("POST", "/test", errReader{})
	wErrReader := httptest.NewRecorder()
	handler.ServeHTTP(wErrReader, reqErrReader)
	if wErrReader.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for manual body error, got %d", wErrReader.Code)
	}
}

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, context.DeadlineExceeded
}
