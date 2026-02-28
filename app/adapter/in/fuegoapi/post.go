package fuegoapi

import (
	"archetype/app/shared/infrastructure/httpserver"

	"github.com/Ignaciojeria/ioc"
	fuegofw "github.com/go-fuego/fuego"
)

var _ = ioc.Register(NewTemplatePost)

type TemplatePostRequest struct {
	Message string `json:"message"`
}

type TemplatePostResponse struct {
	Status string `json:"status"`
}

// NewTemplatePost registers a sample POST endpoint.
func NewTemplatePost(s *httpserver.Server) {
	fuegofw.Post(s.Manager, "/hello",
		func(c fuegofw.ContextWithBody[TemplatePostRequest]) (TemplatePostResponse, error) {
			body, err := c.Body()
			if err != nil {
				return TemplatePostResponse{}, err
			}
			return TemplatePostResponse{
				Status: body.Message + " received",
			}, nil
		})
}
