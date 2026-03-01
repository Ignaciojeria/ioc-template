package pubsub

import (
	"context"
	"errors"
	"fmt"

	"archetype/app/shared/configuration"

	"cloud.google.com/go/pubsub"
	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewClient)

// NewClient creates a new GCP PubSub client using the configuration.
func NewClient(env configuration.PubSubConfiguration) (*pubsub.Client, error) {
	if env.GOOGLE_PROJECT_ID == "" {
		return nil, errors.New("GOOGLE_PROJECT_ID is required for PubSub client")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, env.GOOGLE_PROJECT_ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	return client, nil
}
