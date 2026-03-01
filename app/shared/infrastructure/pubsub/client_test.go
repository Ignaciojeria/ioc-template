package pubsub

import (
	"strings"
	"testing"

	"archetype/app/shared/configuration"
)

func TestNewClient_MissingProjectID(t *testing.T) {
	conf := configuration.Conf{
		GOOGLE_PROJECT_ID: "",
	}

	client, err := NewClient(conf)
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
