package out

import (
	"context"

	"archetype/app/domain/entity"
)

// TemplateRepository defines the persistence contract for Template.
// Implementations (Postgres, MongoDB, etc.) live in adapter/out.
type TemplateRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Template, error)
}
