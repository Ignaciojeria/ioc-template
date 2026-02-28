package fuegoapi

import (
	"archetype/app/shared/infrastructure/httpserver"

	"github.com/Ignaciojeria/ioc"
	fuegofw "github.com/go-fuego/fuego"
)

var _ = ioc.Register(NewTemplatePatch)

type TemplatePatchRequest struct {
	Message string `json:"message"`
}

type TemplatePatchResponse struct {
	Status string `json:"status"`
}

// NewTemplatePatch registers a sample PATCH endpoint.
func NewTemplatePatch(s *httpserver.Server) {
	fuegofw.Patch(s.Manager, "/hello",
		func(c fuegofw.ContextWithBody[TemplatePatchRequest]) (TemplatePatchResponse, error) {
			body, err := c.Body()
			if err != nil {
				return TemplatePatchResponse{}, err
			}
			return TemplatePatchResponse{
				Status: body.Message + " patched",
			}, nil
		})
}
