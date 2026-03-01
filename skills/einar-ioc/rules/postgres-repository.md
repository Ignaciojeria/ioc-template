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

---

## Unit tests

When creating a new component, generate tests following this pattern:

### app/adapter/out/postgres/postgres_repository_test.go

```go
package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestNewTemplateRepository(t *testing.T) {
	repo, err := NewTemplateRepository(nil)
	if err != nil {
		t.Fatalf("expected no error during repository creation, got %v", err)
	}
	if repo == nil {
		t.Fatal("expected repository instance, got nil")
	}
}

func TestTemplateRepository_FindById(t *testing.T) {
	// Create a mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub database connection: %v", err)
	}
	defer db.Close()

	// Wrap the sql.DB into sqlx.DB
	sqlxDB := sqlx.NewDb(db, "postgres")

	// Create repository instance
	repo, _ := NewTemplateRepository(sqlxDB)

	// Subtest 1: Successful fetch
	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM template_table WHERE id = \\$1 LIMIT 1").
			WithArgs("123").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))

		result, err := repo.FindById("123")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result.ID != "123" {
			t.Errorf("expected ID 123, got %s", result.ID)
		}
	})

	// Subtest 2: Fetching fails
	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM template_table WHERE id = \\$1 LIMIT 1").
			WithArgs("999").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.FindById("999")
		if err == nil {
			t.Error("expected error sql.ErrNoRows, got nil")
		}
	})
}
```
