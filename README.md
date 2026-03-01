# Einar Project Template

[![codecov](https://codecov.io/gh/Ignaciojeria/ioc-template/branch/main/graph/badge.svg)](https://codecov.io/gh/Ignaciojeria/ioc-template)
[![Go Report Card](https://goreportcard.com/badge/github.com/Ignaciojeria/ioc-template)](https://goreportcard.com/report/github.com/Ignaciojeria/ioc-template)

Welcome to your new Einar-based Go project! This project has been scaffolded using [Einar CLI](https://github.com/Ignaciojeria/einar) and is powered by [einar-ioc](https://github.com/Ignaciojeria/einar-ioc), a purely Go-idiomatic, type-safe Dependency Injection framework.

This README serves as both a guide for human developers and a context prompt for AI agents collaborating on this codebase.

## 🏗️ Architecture

This project follows a Modular / Hexagonal Architecture (Ports and Adapters) pattern organized by concerns:

- `cmd/api`: The application entry point. Loads dependencies and starts the system.
- `app/shared`: Common infrastructure, configuration, cross-cutting concerns (e.g., HTTP server, OpenTelemetry, logging).
- `app/adapter/in`: Inbound adapters (Controllers, EventBus consumers, GRPC handlers) that trigger business logic.
- `app/adapter/out`: Outbound adapters (Database repositories, HTTP clients, Publishers) to interact with external systems.
- `app/usecase`: Core business logic. Pure Go code, independent of external frameworks.

## 🔌 Inversion of Control (einar-ioc)

This project relies on `einar-ioc` to wire dependencies automatically. 

**Rules for AI Agents and Developers:**
1. **Never use `init()` functions** to instantiate business components.
2. Register components at the package level using `var _ = ioc.Register(...)`.
3. Provide your constructors returning the struct/interface and optionally an `error`.
4. Define your dependencies as function parameters in your constructors.
5. The container will automatically infer and inject the dependencies by matching types (100% Type-Safe).

### Example Component

```go
package mypackage

import "github.com/Ignaciojeria/ioc"

// 1. Register the constructor
var _ = ioc.Register(NewMyService)

type MyService struct {
	db *sqlx.DB // Dependency
}

// 2. Define constructor with dependencies as parameters
func NewMyService(db *sqlx.DB) (*MyService, error) {
	return &MyService{db: db}, nil
}
```

## 🛠️ Einar CLI Usage

You can use the `einar` CLI to incrementally add features and components to your project.

### Installations
Install core infrastructure modules:
- `einar install fuego` (Adds Fuego HTTP Server)
- `einar install postgresql` (Adds PostgreSQL connection via `sqlx` and `golang-migrate`)
- `einar install gcp-pubsub` (Adds Google Cloud PubSub with CloudEvents agnostic adapter)

### Generators
Scaffold standard components (these will be automatically wired into the IoC container):
- `einar generate usecase getUser`
- `einar generate get-controller getUser`
- `einar generate post-controller createUser`
- `einar generate postgres-repository user` (Scaffolds a hexagonal outbound repository with raw SQL)
- `einar generate pubsub-consumer procesar_pedido` (Scaffolds an agnostic CloudEvent listener)

*Note: Einar CLI uses AST parsing to dynamically update your `main.go` with blank imports (`_ "path/to/package"`) when generating or installing components. You rarely need to modify `cmd/api/main.go` manually.*

## 🧪 Testing

The IoC container is designed to facilitate high unit test coverage. Since constructors explicitly define dependencies as parameters, you can easily instantiate them in tests by passing mocks or stubs directly.

**Avoid extracting variables for testability:**
Do **not** extract package-level variables (e.g. `var jsonMarshal = json.Marshal`) or replace standard library calls with injectable references solely to reach 100% coverage. This practice:
- Introduces global mutable state that can be overridden in tests
- Hides dependencies and pollutes production code
- Is considered a code smell in idiomatic Go

Prefer constructor injection or interfaces for dependency injection. Accept slightly lower coverage (e.g. 97–98%) for error paths that are practically unreachable rather than degrading code quality.

**Integration Testing:**
This template incorporates `testcontainers-go`. When running PostgreSQL tests, a real Docker container is automatically spun up to validate that:
1. SQL migrations are applied correctly.
2. Raw SQL queries work against a real engine.

```bash
# Run all tests (including integration tests with Docker)
go test -v ./...
```

## 🗄️ Database & Migrations

We prioritize **Raw SQL** and **Explicit Migrations** over ORM magic to ensure 100% control over the schema and performance.

**Rule for AI Agents & Developers:**
1. **Never use ORMs with Auto-Migrations.** 
2. **Schema Management:** All table changes MUST reside in pure `.sql` files within the `app/shared/infrastructure/postgresql/migrations/` directory.
3. **Execution:** Migrations are embedded into the Go binary using `//go:embed` and executed on application startup using `golang-migrate`.
4. **Data Access:** Use `sqlx` for data access. Map rows directly to structs using `db:"column_name"` tags.
5. **Testing:** Use `go-sqlmock` for unit testing repositories and `testcontainers` for integration testing.

## 📡 Messaging & PubSub (CloudEvents)

To avoid vendor lock-in, all messaging components follow the **CloudEvents (CNCF)** standard.

- **Subscriber Strategy**: Uses an agnostic `eventbus.Subscriber` interface. 
- **Adapters**: The GCP PubSub adapter handles PULL from the cloud and provides a PUSH fallback via HTTP (WebHooks) for local debugging.
- **Payloads**: All consumers receive `cloudevents.Event`, keeping the business logic independent of the broker (PubSub, NATS, Kafka).

## 📄 Environment Configuration

The application manages variables via a single connection string or specific envs. Configuration is processed in `app/shared/configuration/`.

- Default configurations use Struct Tags (e.g., `` env:"DATABASE_URL" envDefault:"postgres://..." ``).
- Use `configuration.Conf` as a dependency. Never read from `os.Getenv` directly within logic.

## 📝 Generated Documentation

The `skills/einar-ioc/rules/` directory contains markdown files with the full source code of each template component. These files are auto-generated from the actual codebase and can be used as Cursor rules or skills.

**Regenerate after modifying template files:**

```bash
go run ./scripts/gen-from-template
```

The mapping between source files and generated markdown is defined in `scripts/gen-skills.config.yaml`. To add a new component, add an entry to the config and re-run the script.
