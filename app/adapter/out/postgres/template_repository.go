package postgres

import (
	"context"

	"archetype/app/application/ports/out"
	"archetype/app/domain/entity"

	"github.com/Ignaciojeria/ioc"
	"github.com/jmoiron/sqlx"
)

var _ = ioc.Register(NewTemplateRepository)

type templateRepository struct {
	db *sqlx.DB
}

// NewTemplateRepository returns an implementation of ports/out.TemplateRepository.
func NewTemplateRepository(db *sqlx.DB) (out.TemplateRepository, error) {
	return &templateRepository{db: db}, nil
}

func (r *templateRepository) FindByID(ctx context.Context, id string) (*entity.Template, error) {
	var dest struct {
		ID string `db:"id"`
	}
	err := r.db.GetContext(ctx, &dest, "SELECT id FROM template_table WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	return &entity.Template{ID: dest.ID}, nil
}
