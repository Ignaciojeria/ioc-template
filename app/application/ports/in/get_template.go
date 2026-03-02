package in

import "context"

// GetTemplateExecutor is the input port for the GetTemplate use case.
// Controllers (Fuego, gRPC) call this; implementations live in usecase/.
type GetTemplateExecutor interface {
	Execute(ctx context.Context, id string) (GetTemplateOutput, error)
}

// GetTemplateOutput is the DTO returned by the use case.
type GetTemplateOutput struct {
	ID   string
	Name string
}
