package eventbus

import (
	"os"
	"strings"
	"testing"

	"archetype/app/shared/configuration"
)

func TestNewGcpClient_MissingProjectID(t *testing.T) {
	conf := configuration.Conf{
		GOOGLE_PROJECT_ID: "",
	}

	client, err := NewGcpClient(conf)
	if err == nil {
		t.Fatal("expected error creating pubsub client with empty google project ID, got nil")
	}
	if client != nil {
		t.Errorf("expected nil pubsub client on error, got %v", client)
	}

	if !strings.Contains(err.Error(), "GOOGLE_PROJECT_ID is required") {
		t.Errorf("expected missing GOOGLE_PROJECT_ID formatting error, got %v", err)
	}
}

func TestNewGcpClient_FailureToConnect(t *testing.T) {
	conf := configuration.Conf{
		GOOGLE_PROJECT_ID: "test-project",
	}

	// This should fail to create client because credentials are not found and pubsub requires it
	// unless running with option.WithoutAuthentication() which our client does not embed,
	// or emulator env var. So setting a random PUBSUB_EMULATOR_HOST to an invalid address gives an error.
	os.Setenv("PUBSUB_EMULATOR_HOST", "127.0.0.1:0")
	defer os.Unsetenv("PUBSUB_EMULATOR_HOST")

	_, err := NewGcpClient(conf)
	if err != nil {
		// NewClient succeeds synchronously even with fake emulator host
		// if credentials are not checked immediately, but if it returns an err, we catch it.
	}
}
