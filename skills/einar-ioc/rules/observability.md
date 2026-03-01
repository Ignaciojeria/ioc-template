---
name: observability
description: Distributed tracing and JSON logging via Dependency Injection
---

## Overview

Einar provides unified observability through OpenTelemetry (traces) and `log/slog` (structured JSON logging).

These components are encapsulated inside the `observability.Observability` struct defined in `app/shared/infrastructure/observability`.

## Injecting Observability

You MUST NEVER use the standard `log` package, `fmt.Printf`, or `fmt.Println` to log information.
Instead, your components MUST inject `observability.Observability` via their IoC constructors.

### Boilerplate Example

```go
package usecase

import (
    "context"
    "github.com/Ignaciojeria/ioc"
    "archetype/app/shared/infrastructure/observability"
)

var _ = ioc.Register(NewMyUseCase)

type MyUseCase struct {
    obs observability.Observability
}

// 1. Ask for Observability in the constructor
func NewMyUseCase(obs observability.Observability) *MyUseCase {
    return &MyUseCase{obs: obs}
}
```

## Context-Aware Logging (Slog)

To ensure logs are correlated with OpenTelemetry traces (TraceID, SpanID), you MUST use the contextual logging methods (`InfoContext`, `ErrorContext`, etc.) on the injected logger.

```go
func (uc *MyUseCase) Execute(ctx context.Context) {
    // CORRECT: The context contains the TraceID injected by Fuego or EventBus middlewares
    uc.obs.Logger.InfoContext(ctx, "started executing use case",
        "action", "Execute",
        "attempt", 1,
    )

    // WRONG: Loses trace correlation!
    uc.obs.Logger.Info("started executing use case")
}
```

## OpenTelemetry Tracing

When executing complex business logic or wrapping external calls, you MUST use the injected `Tracer` to create spans natively.

```go
import "go.opentelemetry.io/otel/attribute"

func (uc *MyUseCase) Execute(ctx context.Context, param string) error {
    // 1. Start a span from the injected Tracer
    ctx, span := uc.obs.Tracer.Start(ctx, "ExecuteLogic")
    defer span.End() // 2. Ensure it ends

    // 3. Add custom attributes for searchability
    span.SetAttributes(attribute.String("param.value", param))

    // 4. Pass the wrapped context down to repositories or publishers 
    // to keep the distributed trace alive natively!
    err := uc.repository.Save(ctx, param)
    if err != nil {
        span.RecordError(err)
        uc.obs.Logger.ErrorContext(ctx, "failed to save", "error", err)
        return err
    }

    return nil
}
```
