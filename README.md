# Einar Project Template

Welcome to your new Einar-based Go project! This project has been scaffolded using [Einar CLI](https://github.com/Ignaciojeria/einar) and is powered by [einar-ioc](https://github.com/Ignaciojeria/einar-ioc), a purely Go-idiomatic, type-safe Dependency Injection framework.

This README serves as both a guide for human developers and a context prompt for AI agents collaborating on this codebase.

## 🏗️ Architecture

This project follows a Modular / Hexagonal Architecture (Ports and Adapters) pattern organized by concerns:

- `cmd/api`: The application entry point. Loads dependencies and starts the system.
- `app/shared`: Common infrastructure, configuration, cross-cutting concerns (e.g., HTTP server, OpenTelemetry, logging).
- `app/adapter/in`: Inbound adapters (Controllers, PubSub consumers, GRPC handlers) that trigger business logic.
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
	db Database // Dependency
}

// 2. Define constructor with dependencies as parameters
func NewMyService(db Database) (*MyService, error) {
	return &MyService{db: db}, nil
}
```

## 🛠️ Einar CLI Usage

You can use the `einar` CLI to incrementally add features and components to your project.

### Installations
Install core infrastructure modules:
- `einar install fuego` (Adds Fuego HTTP Server)
- `einar install postgresql` (Adds GORM PostgreSQL connection)
- `einar install gcp-pubsub` (Adds Google Cloud PubSub)

### Generators
Scaffold standard components (these will be automatically wired into the IoC container):
- `einar generate usecase getUser`
- `einar generate get-controller getUser`
- `einar generate post-controller createUser`

*Note: Einar CLI uses AST parsing to dynamically update your `main.go` with blank imports (`_ "path/to/package"`) when generating or installing components. You do rarely need to modify `main.go` manually.*

## 🧪 Testing

The IoC container is designed to facilitate 100% unit test coverage. Since constructors explicitly define dependencies as parameters, you can easily instantiate them in tests by passing mocks or stubs directly—without needing to boot the full IoC container.

```go
func TestNewMyService(t *testing.T) {
	mockDB := NewMockDatabase()
	service, err := NewMyService(mockDB)
    // assert...
}
```

## 📄 Environment Configuration

The application expects variables to be managed via `.env` files for local development. Configuration is processed in `app/shared/configuration/conf.go`.

- Default configurations use Struct Tags (e.g., `` env:"PORT" envDefault:"8080" ``).
- Use `configuration.NewConf` as a dependency if your component needs to access environment variables. Never read from `os.Getenv` directly within business logic.
