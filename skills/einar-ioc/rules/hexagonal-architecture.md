---
name: hexagonal-architecture
description: Project structure and architectural limits
---

## Directory Structure

When working with an Einar-IoC project, you MUST strictly adhere to this Hexagonal Architecture directory layout:

```text
app/
├── adapter/
│   ├── in/         # Incoming adapters (HTTP controllers, event consumers, CLI commands)
│   └── out/        # Outgoing adapters (Database repositories, third-party API clients, event publishers)
├── core/
│   ├── usecase/    # Application business logic (interactors coordinating domain and out ports)
│   └── domain/     # Core domain entities and pure business rules (no external dependencies)
└── shared/         # Cross-cutting concerns (configuration, logging, tracing, utilities)
```

## Architectural Constraints (CRITICAL)

1. **Dependency Inversion**: 
   - `core` (`domain` and `usecase`) MUST NOT import anything from `adapter` (`in` or `out`). It defines interfaces (Ports) that adapters must implement.
   - `adapter/in` MUST NOT import anything from `adapter/out`. Incoming requests must flow through `core/usecase`.
   - `adapter/out` implements interfaces defined in `core` and connects to external systems.

2. **Incoming Adapters (`app/adapter/in`)**:
   - Includes Fuego HTTP controllers (`app/adapter/in/fuegoapi/...`).
   - Includes Event Subscribers (`app/adapter/in/eventbus/...`).
   - **Responsibility**: Parse incoming requests/events, call a `core/usecase`, and return a response. They DO NOT contain business logic nor interact directly with databases.

3. **Outgoing Adapters (`app/adapter/out`)**:
   - Includes PostgreSQL/GORM repositories (`app/adapter/out/postgres/...`).
   - Includes Event Publishers (`app/adapter/out/eventbus/...`).
   - Includes external REST API clients.
   - **Responsibility**: Implement interfaces defined by `core` to fetch/mutate data in external systems.

4. **Shared Concerns (`app/shared`)**:
   - Includes configuration loading (`app/shared/configuration`).
   - Includes observability logic like slog loggers and trace propagation (`app/shared/infrastructure/observability`).
   - Can be imported by ANY layer.

**Failure to respect these boundaries will result in spaghetti code. Enforce Hexagonal constraints strictly.**
