package postgres

import (
	"github.com/Ignaciojeria/ioc"
	"github.com/jmoiron/sqlx"
)

var _ = ioc.Register(NewTemplateRepository)

// TemplateStructRepository fetches TemplateStruct by ID. Implemented by *TemplateRepository.
type TemplateStructRepository interface {
	FindById(id string) (TemplateStruct, error)
}

type TemplateStruct struct {
	ID string `db:"id"`
}

type TemplateRepository struct {
	db *sqlx.DB
}

func NewTemplateRepository(db *sqlx.DB) (TemplateStructRepository, error) {
	return &TemplateRepository{db: db}, nil
}

// FindById fetches a TemplateStruct by its ID.
func (r *TemplateRepository) FindById(id string) (TemplateStruct, error) {
	var dest TemplateStruct
	// sqlx allows mapping raw queries into structs easily using db tags.
	err := r.db.Get(&dest, "SELECT * FROM template_table WHERE id = $1 LIMIT 1", id)
	return dest, err
}
