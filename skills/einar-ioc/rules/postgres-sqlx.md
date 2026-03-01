---
name: postgres-sqlx
description: Database access rules using PostgreSQL and sqlx
---

## Overview

Einar uses PostgreSQL with `github.com/jmoiron/sqlx` for database access.
All database interaction MUST occur inside `app/adapter/out/postgres`. 
You MUST NOT place database logic inside controllers, event subscribers, or usecases.

## The Repository Template

When implementing a database repository, inject `*sqlx.DB` via IoC. 

### Boilerplate

```go
package postgres

import (
	"context"
	
	"github.com/Ignaciojeria/ioc"
	"github.com/jmoiron/sqlx"
)

var _ = ioc.Register(NewUserRepository)

// Private to the postgres adapter to prevent domain leakage
type userModel struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func (m userModel) ToDomain() domain.User {
	return domain.User{
		ID:   m.ID,
		Name: m.Name,
	}
}

type UserRepository struct {
	db *sqlx.DB
}

// The constructor receives the global sqlx DB instance
func NewUserRepository(db *sqlx.DB) (*UserRepository, error) {
	return &UserRepository{db: db}, nil
}

// Method fulfilling the outbound port interface defined in `core`
func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	var user userModel
	
	// Use sqlx's GetContext to map the raw query securely into the struct, passing Context.
	// Context is required so OpenTelemetry can properly trace DB query spans.
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return domain.User{}, err
	}
	
	return user.ToDomain(), nil
}
```

## SQLx Methods

You MUST use SQLx's context-aware methods for database queries to ensure Fuego's OpenTelemetry traces map properly to database queries:
- **`GetContext`**: Use when expecting exactly a single row.
- **`SelectContext`**: Use when expecting multiple rows.
- **`ExecContext`**: Use for `INSERT`, `UPDATE`, `DELETE`.
- **`QueryRowxContext` / `QueryxContext`**: Use when more control is needed.

## SQLx Models vs Domain Entities

We strictly separate database structs (SQLx models) from domain entities.
Your SQLx models should be unexported (or kept private to the `postgres` package) to prevent leakage. Use the `db` struct tags to automatically map Postgres columns directly to your struct fields.
