---
name: einar-ioc
description: Best practices for building Go applications using Hexagonal Architecture and Einar-IoC
metadata:
  tags: go, golang, hexagonal-architecture, dependency-injection, ioc, fuego
---

## When to use

Use this skill whenever you are writing or modifying Go code that uses the `github.com/Ignaciojeria/ioc` library, or when you are creating APIs, event-driven components, or database adapters in a Hexagonal Architecture.

> [!TIP]
> **Living Documentation**: Since this skill lives inside the `ioc-template` repository (or a codebase spawned from it), you can read the actual `.go` files inside `app/adapter` and `app/shared` to see exactly how these rules are applied in practice as your primary reference!

## Architecture Guidelines

Strictly follow the Hexagonal Architecture constraints. Read [./rules/hexagonal-architecture.md](./rules/hexagonal-architecture.md) for directory structures and import limitations.

## Dependency Injection (IoC)

All components MUST be registered in the IoC container instead of using manual wiring in `main.go`. Read [./rules/dependency-injection.md](./rules/dependency-injection.md) for the exact `var _ = ioc.Register(...)` syntax.

## Building HTTP APIs (Fuego)

Whenever you need to expose a REST endpoint, you MUST use `github.com/go-fuego/fuego`. Do not use standard `net/http`, `gin`, or `fiber`. Read [./rules/fuego-controllers.md](./rules/fuego-controllers.md) for the controller template.

## Event-Driven Messaging (EventBus)

When dealing with asynchronous messaging (Pub/Sub, queues), implement the standardized EventBus interfaces using CloudEvents. Read [./rules/eventbus-messaging.md](./rules/eventbus-messaging.md) to understand how `Publisher` and `Subscriber` adapters work under the Multi-Broker Factory (GCP/NATS).

## Database Adapters (PostgreSQL)

When creating database repositories, use SQLx and inject the `*sqlx.DB` instance. Read [./rules/postgres-sqlx.md](./rules/postgres-sqlx.md).

## Observability (OpenTelemetry & Slog)

All components must log via an injected `*slog.Logger` from the `shared` observability package. HTTP calls and messaging must propagate OpenTelemetry Context. Read [./rules/observability.md](./rules/observability.md).

## How to use

Read individual rule files for detailed explanations and code examples:

- [rules/hexagonal-architecture.md](rules/hexagonal-architecture.md) - Project structure and architectural limits.
- [rules/dependency-injection.md](rules/dependency-injection.md) - How to use `einar-ioc` for automatic wiring.
- [rules/fuego-controllers.md](rules/fuego-controllers.md) - Boilerplate for defining REST APIs.
- [rules/eventbus-messaging.md](rules/eventbus-messaging.md) - Handling GCP and NATS messaging via CloudEvents.
- [rules/postgres-sqlx.md](rules/postgres-sqlx.md) - Database access rules.
- [rules/observability.md](rules/observability.md) - Distributed tracing and JSON logging.
