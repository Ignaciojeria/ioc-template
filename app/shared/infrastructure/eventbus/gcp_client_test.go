package eventbus

import (
	"context"
	"errors"
	"strings"
	"testing"

	"archetype/app/shared/configuration"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func TestNewGcpClient_MissingProjectID(t *testing.T) {
	conf := configuration.Conf{GOOGLE_PROJECT_ID: ""}

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

func TestNewGcpClient_CreateClientError(t *testing.T) {
	orig := pubsubNewClient
	defer func() { pubsubNewClient = orig }()

	pubsubNewClient = func(ctx context.Context, projectID string, opts ...option.ClientOption) (*pubsub.Client, error) {
		return nil, errors.New("new client fail")
	}

	_, err := NewGcpClient(configuration.Conf{GOOGLE_PROJECT_ID: "project"})
	if err == nil {
		t.Fatal("expected create client error")
	}
}

func TestNewGcpClient_Success(t *testing.T) {
	orig := pubsubNewClient
	defer func() { pubsubNewClient = orig }()

	dummy := &pubsub.Client{}
	pubsubNewClient = func(ctx context.Context, projectID string, opts ...option.ClientOption) (*pubsub.Client, error) {
		return dummy, nil
	}

	client, err := NewGcpClient(configuration.Conf{GOOGLE_PROJECT_ID: "project"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client != dummy {
		t.Fatal("expected mocked client instance")
	}
}
