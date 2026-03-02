package usecase

import (
	"context"

	"archetype/app/adapter/out/postgres"

	"github.com/Ignaciojeria/ioc"
)

var _ = ioc.Register(NewTemplateUseCase)

// TemplateExecutor defines the contract for the template use case.
// Implementations are injected into controllers via IoC.
type TemplateExecutor interface {
	Execute(ctx context.Context, id string) (TemplateOutput, error)
}

// TemplateOutput is the output DTO for the template use case.
// Use case-specific DTOs (XxxInput, XxxOutput) are preferred over domain entities because:
// they decouple the API from the domain, each use case exposes only what it needs, and
// they are easier for agents to generate (self-contained file, predictable pattern).
type TemplateOutput struct {
	ID   string
	Name string
}

type templateUseCase struct {
	repo postgres.TemplateStructRepository
}

func NewTemplateUseCase(repo postgres.TemplateStructRepository) (TemplateExecutor, error) {
	return &templateUseCase{repo: repo}, nil
}

func (uc *templateUseCase) Execute(ctx context.Context, id string) (TemplateOutput, error) {
	t, err := uc.repo.FindById(id)
	if err != nil {
		return TemplateOutput{}, err
	}
	return TemplateOutput{ID: t.ID}, nil
}
