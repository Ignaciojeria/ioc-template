package eventbus

import (
	"context"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func TestApplyMiddlewares(t *testing.T) {
	// Base processor
	baseProcessor := func(ctx context.Context, event cloudevents.Event) int {
		return 200
	}

	// Middleware 1
	m1 := func(next MessageProcessor) MessageProcessor {
		return func(ctx context.Context, event cloudevents.Event) int {
			return next(ctx, event) + 10
		}
	}

	// Middleware 2
	m2 := func(next MessageProcessor) MessageProcessor {
		return func(ctx context.Context, event cloudevents.Event) int {
			return next(ctx, event) + 5
		}
	}

	// Wrapped
	final := ApplyMiddlewares(baseProcessor, m1, m2)

	ce := cloudevents.NewEvent()
	res := final(context.Background(), ce)

	// Since m1 and m2 add 10 and 5, the total should be 215
	if res != 215 {
		t.Errorf("expected 215, got %d", res)
	}
}
