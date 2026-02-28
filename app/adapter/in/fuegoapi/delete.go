package fuegoapi

import (
	"archetype/app/shared/infrastructure/httpserver"

	"github.com/Ignaciojeria/ioc"
	fuegofw "github.com/go-fuego/fuego"
)

var _ = ioc.Register(NewTemplateDelete)

// NewTemplateDelete registers a sample DELETE endpoint.
func NewTemplateDelete(s *httpserver.Server) {
	fuegofw.Delete(s.Manager, "/hello/{id}",
		func(c fuegofw.ContextNoBody) (string, error) {
			id := c.PathParam("id")
			return id + " deleted", nil
		})
}
