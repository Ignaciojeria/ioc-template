# architecture-guidelines

> Correct patterns for interfaces, ports, and injection with einar-ioc

## No ports folder

**Do not create `app/domain/port/` or similar folders.** Interfaces should live next to their implementation or in the consumer that defines them.

- **Interface in the consumer (idiomatic Go):** The usecase defines what it needs: `type AccountRepository interface { GetByID(...) }` in the same package as the usecase, or in an adjacent file.
- **Interface next to the implementation:** In `app/adapter/out/postgres/account_repository.go` you can define the interface that adapter implements if it is the only consumer.

Example: **repository** (interface next to implementation) and **usecase** (defines what it needs, returns its interface):

```go
// app/adapter/out/postgres/account_repository.go
package postgres

import (
	"context"
	"github.com/Ignaciojeria/ioc"
	"github.com/jmoiron/sqlx"
)

var _ = ioc.Register(NewAccountRepository)

type AccountRepository interface {
	GetByID(ctx context.Context, id string) (*Account, error)
}

type Account struct {
	ID string `db:"id"`
}

type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) (AccountRepository, error) {
	return &accountRepository{db: db}, nil
}

func (r *accountRepository) GetByID(ctx context.Context, id string) (*Account, error) {
	var dest Account
	err := r.db.GetContext(ctx, &dest, "SELECT * FROM accounts WHERE id = $1", id)
	return &dest, err
}
```

```go
// app/application/usecase/get_account.go
package usecase

import (
	"context"
	"archetype/app/adapter/out/postgres"
	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewGetAccountUseCase)

type GetAccountExecutor interface {
	Execute(ctx context.Context, id string) (*AccountOutput, error)
}

type getAccountUseCase struct {
	repo postgres.AccountRepository
}

func NewGetAccountUseCase(repo postgres.AccountRepository) (GetAccountExecutor, error) {
	return &getAccountUseCase{repo: repo}, nil
}
```

The postgres adapter owns `AccountRepository` (interface + implementation). The usecase injects it and exposes `GetAccountExecutor`. The controller injects `GetAccountExecutor`. No separate `port` package is needed.

## Always inject interfaces (one implementation per interface)

**Always inject and return interfaces** in IoC constructors. Controllers, use cases, and adapters depend on contracts, not concrete types. The constructor returns the interface it implements so consumers depend on the abstraction.

**To avoid ambiguity:** Each interface must have **exactly one implementation** registered in the IoC. If an interface has N implementations, the container cannot know which one to inject.

```go
// app/adapter/in/fuegoapi/get_account.go
// ✅ Controller: inject interface (same example as above)
func NewGetAccount(s *httpserver.Server, uc usecase.GetAccountExecutor) {
	fuegofw.Get(s.Manager, "/accounts/{id}",
		func(c fuegofw.ContextNoBody) (GetAccountResponse, error) {
			out, err := uc.Execute(c.Context(), c.PathParam("id"))
			// ...
			return GetAccountResponse{ID: out.ID}, nil
		},
	)
}
```

Repository and use case are shown in the example above. Full flow: `AccountRepository` → `GetAccountUseCase` → `GetAccount` controller.

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
| Injecting or returning concrete types in constructors | Always inject and return interfaces |
| One interface with N implementations (ambiguity) | One interface → one implementation; use factory or separate interfaces if needed |
| Returning domain entities from use cases | Use case-specific DTOs (`XxxInput`, `XxxOutput`) |
