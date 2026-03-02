# architecture-guidelines

> Correct patterns for interfaces, ports, and injection with einar-ioc

## No ports folder

**Do not create `app/domain/port/` or similar folders.** Interfaces should live next to their implementation or in the consumer that defines them.

- **Interface in the consumer (idiomatic Go):** The usecase defines what it needs: `type AccountRepository interface { GetByID(...) }` in the same package as the usecase, or in an adjacent file.
- **Interface next to the implementation:** In `app/adapter/out/postgres/account_repository.go` you can define the interface that adapter implements if it is the only consumer.

Example: the usecase defines the contract it needs:

```go
// app/application/usecase/get_account.go
package usecase

type AccountRepository interface {
	GetByID(ctx context.Context, id string) (*Account, error)
}

type GetAccountUseCase struct {
	repo AccountRepository
}

func NewGetAccountUseCase(repo AccountRepository) (*GetAccountUseCase, error) {
	return &GetAccountUseCase{repo: repo}, nil
}
```

The postgres adapter implements that interface. No separate `port` package is needed.

## Always inject interfaces (one implementation per interface)

**Always inject the interface** in IoC constructors. Controllers, use cases, and adapters depend on contracts, not concrete types.

**To avoid ambiguity:** Each interface must have **exactly one implementation** registered in the IoC. If an interface has N implementations, the container cannot know which one to inject.

```go
// ✅ Controller: inject interface
func NewGetTemplate(s *httpserver.Server, uc usecase.TemplateExecutor) {
	fuegofw.Get(s.Manager, "/templates/{id}",
		func(c fuegofw.ContextNoBody) (GetTemplateResponse, error) {
			out, err := uc.Execute(c.Context(), c.PathParam("id"))
			// ...
		},
	)
}

// ✅ Use case: inject interface (repository)
func NewTemplateUseCase(repo postgres.TemplateStructRepository) (usecase.TemplateExecutor, error) {
	return &templateUseCase{repo: repo}, nil
}
```

If you need multiple implementations of the same contract (e.g. different event publishers), use separate interfaces or a factory pattern.

## Use case output: DTOs over domain entities

**Prefer use case-specific DTOs** (`XxxInput`, `XxxOutput`) over returning domain entities:

- **Decoupling:** The API contract stays independent of the domain model.
- **Explicit contracts:** Each use case exposes only the fields it needs.
- **Agent-friendly:** DTOs allow self-contained use case files with predictable patterns; agents can generate complete use cases without cross-package lookups.

See `app/application/usecase/interfaces.go` for the pattern.

## Summary

| Avoid | Use |
|-------|-----|
| `app/domain/port/` with standalone interfaces | Interface in the consumer or next to the implementation |
| Injecting concrete types in constructors | Always inject interfaces |
| One interface with N implementations (ambiguity) | One interface → one implementation; use factory or separate interfaces if needed |
| Returning domain entities from use cases | Use case-specific DTOs (`XxxInput`, `XxxOutput`) |
