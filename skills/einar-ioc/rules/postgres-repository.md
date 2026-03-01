# postgres-repository

> PostgreSQL repository pattern with sqlx and go-sqlmock

## app/adapter/out/postgres/postgres_repository.go

```go
package postgres

import (
	"github.com/Ignaciojeria/ioc"
	"github.com/jmoiron/sqlx"
)

var _ = ioc.Register(NewTemplateRepository)

type TemplateStruct struct {
	ID string `db:"id"`
}

type TemplateRepository struct {
	db *sqlx.DB
}

func NewTemplateRepository(db *sqlx.DB) (*TemplateRepository, error) {
	return &TemplateRepository{db: db}, nil
}

// FindById fetches a TemplateStruct by its ID.
func (r *TemplateRepository) FindById(id string) (TemplateStruct, error) {
	var dest TemplateStruct
	// sqlx allows mapping raw queries into structs easily using db tags.
	err := r.db.Get(&dest, "SELECT * FROM template_table WHERE id = $1 LIMIT 1", id)
	return dest, err
}
```
