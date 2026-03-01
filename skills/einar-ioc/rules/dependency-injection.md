---
name: dependency-injection
description: How to use Einar-IoC for automatic wiring
---

## Overview

Einar uses the `github.com/Ignaciojeria/ioc` library to provide native Dependency Injection. 
You MUST NOT wire dependencies manually (e.g., passing nested dependencies down a chain in `main.go`).

## How to Register Components

Every struct, client, controller, or repository you create MUST be registered in the IoC container. To do this, wrap your constructor function with `ioc.Register()` at the package level.

### Example: Registering a standard struct

```go
package customer

import "github.com/Ignaciojeria/ioc"

// 1. Declare the registry at the top of the file
var _ = ioc.Register(NewCustomerService)

type CustomerService struct {
	repo Repository // Injected interface
}

// 2. The constructor must return the struct (or interface) and optionally an error
func NewCustomerService(repo Repository) (*CustomerService, error) {
	return &CustomerService{repo: repo}, nil
}
```

### Constructor Rules

1. **Parameters**: The arguments to your constructor function (e.g., `repo Repository`) represent the dependencies your component needs. The IoC container will automatically find and inject matching implementations.
2. **Return values**: Constructors must return `(T, error)` or just `T`, where `T` is a pointer to a struct or an interface. Do not return concrete value types (e.g., return `*MyStruct` not `MyStruct`).

## Interfaces vs Implementations

If your constructor requires an interface, the IoC container will inject the struct that implements that interface.

For example, if `core/usecase` requires `eventbus.Publisher`:
```go
// Inside core/usecase/publish_message.go
var _ = ioc.Register(NewPublishMessageUseCase)

func NewPublishMessageUseCase(pub eventbus.Publisher) *PublishMessageUseCase { ... }
```
And the adapter implements it:
```go
// Inside adapter/out/eventbus/gcp_publisher.go
var _ = ioc.Register(NewGcpPublisher)

// As long as *GcpPublisher implements eventbus.Publisher, it will be injected automatically!
func NewGcpPublisher(client *pubsub.Client) (*GcpPublisher, error) { ... }
```

**CRITICAL RULE:** NEVER create `main.go` wires manually. ALWAYS use `var _ = ioc.Register(NewFunction)` and let the IoC initialize everything in the background.
