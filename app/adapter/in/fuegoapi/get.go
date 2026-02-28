package fuegoapi

import (
	"archetype/app/shared/infrastructure/httpserver"

	"github.com/Ignaciojeria/ioc"
	fuegofw "github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

var _ = ioc.Register(NewTemplateGet)

// NewTemplateGet registers a sample GET endpoint.
// It uses *fuego.Server as a dependency, ensuring the server is allocated before this runs.
func NewTemplateGet(s *httpserver.Server) {
	fuegofw.Get(s.Manager, "/hello",
		func(c fuegofw.ContextNoBody) (string, error) {
			return "Hello from Einar IoC Fuego Template!", nil
		},
		option.Summary("newTemplateGet"),
	)
}
