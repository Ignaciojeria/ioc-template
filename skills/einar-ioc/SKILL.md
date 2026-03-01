---
name: einar-ioc
description: Best practices for building Go applications using Hexagonal Architecture and Einar-IoC
metadata:
  tags: go, golang, hexagonal-architecture, dependency-injection, ioc, fuego, einar-cli, vibe-coding
---

## Purpose

This skill enables **vibe-coding** with the einar-ioc template: the AI follows these rules to scaffold, wire, and modify code correctly without requiring the developer to run Einar CLI. Use it when working on projects based on this template.

## When to use

Use this skill whenever you are writing or modifying Go code that uses the `github.com/Ignaciojeria/ioc` library, or when you are creating APIs, event-driven components, or database adapters in a Hexagonal Architecture. When adding new components (controllers, repositories, consumers, etc.), apply the patterns from the rules—including adding blank imports to `cmd/api/main.go`—as if Einar CLI had scaffolded them.

> [!TIP]
> **Living Documentation**: The rules in `./rules/` are generated from the actual `.go` files. Read the source code in `app/adapter` and `app/shared` as your primary reference—the rules reflect the template as it exists.

## Quick reference

| Domain           | Rule files |
|------------------|------------|
| Structure & main | [structure](rules/structure.md), [main](rules/main.md), [archetype-version](rules/archetype-version.md) |
| Configuration    | [configuration](rules/configuration.md) |
| HTTP / REST      | [httpserver](rules/httpserver.md), [request-logger-middleware](rules/request-logger-middleware.md), [fuegoapi-controllers](rules/fuegoapi-controllers.md) |
| EventBus         | [eventbus-strategy](rules/eventbus-strategy.md), [eventbus-gcp](rules/eventbus-gcp.md), [eventbus-nats](rules/eventbus-nats.md), [consumer](rules/consumer.md), [publisher](rules/publisher.md) |
| Database         | [postgresql-connection](rules/postgresql-connection.md), [postgresql-migrations](rules/postgresql-migrations.md), [postgres-repository](rules/postgres-repository.md) |
| Observability    | [observability](rules/observability.md) |

## Dependency Injection (IoC)

All components MUST be registered in the IoC container via `var _ = ioc.Register(Constructor)`. See any rule file (e.g. [httpserver](rules/httpserver.md), [fuegoapi-controllers](rules/fuegoapi-controllers.md)) for the exact pattern.

**Blank imports are critical:** Each package that registers constructors MUST be imported in `cmd/api/main.go` via a blank import (`_ "archetype/path/to/package"`). Without it, the package never loads and the IoC container will not receive those constructors. When adding a new component, add the corresponding blank import.

## Constraints (from README)

1. **No `init()`** for business components. Use `ioc.Register` at package level instead.
2. **No `os.Getenv` in logic.** Inject `configuration.Conf` as a dependency.
3. **No ORMs with auto-migrations.** Schema changes go in `.sql` files under `app/shared/infrastructure/postgresql/migrations/`. Use `sqlx` + `golang-migrate`.
4. **Don't extract variables for testability.** Avoid `var jsonMarshal = json.Marshal` or injectable stubs just to reach 100% coverage. Prefer constructor injection; accept slightly lower coverage for unreachable error paths.

## Rules by domain

### Structure and entry point

- [structure](rules/structure.md) – Project directory tree (`app/adapter`, `app/shared`, `cmd`, `scripts`)
- [main](rules/main.md) – Entry point and `ioc.LoadDependencies()`
- [archetype-version](rules/archetype-version.md) – Embedded `Version` from `.version` file

### Configuration

- [configuration](rules/configuration.md) – Environment config with caarlos0/env and godotenv

### HTTP / REST (Fuego)

Use `github.com/go-fuego/fuego` for REST. Do not use `net/http`, `gin`, or `fiber` directly.

- [httpserver](rules/httpserver.md) – Fuego server, healthcheck, graceful shutdown
- [request-logger-middleware](rules/request-logger-middleware.md) – HTTP request logging
- [fuegoapi-controllers](rules/fuegoapi-controllers.md) – REST controller scaffold (GET, POST, PUT, PATCH, DELETE)

### EventBus (CloudEvents, GCP / NATS)

- [eventbus-strategy](rules/eventbus-strategy.md) – Interfaces, factory, and strategy pattern
- [eventbus-gcp](rules/eventbus-gcp.md) – GCP Pub/Sub client, publisher, subscriber
- [eventbus-nats](rules/eventbus-nats.md) – NATS client, publisher, subscriber
- [consumer](rules/consumer.md) – Inbound consumer pattern
- [publisher](rules/publisher.md) – Outbound publisher adapter

### Database (PostgreSQL)

- [postgresql-connection](rules/postgresql-connection.md) – SQLx + golang-migrate, connection lifecycle
- [postgresql-migrations](rules/postgresql-migrations.md) – Example migrations: naming (`NNNNNN_name.up/down.sql`), up/down pattern
- [postgres-repository](rules/postgres-repository.md) – Repository pattern with sqlx and go-sqlmock

### Observability

- [observability](rules/observability.md) – OpenTelemetry, slog injection, context propagation
