package fuegoapi

import (
	"archetype/app/shared/infrastructure/httpserver"

	"github.com/Ignaciojeria/ioc"
	fuegofw "github.com/go-fuego/fuego"
)

var _ = ioc.Register(NewTemplatePut)

type TemplatePutRequest struct {
	Message string `json:"message"`
}

type TemplatePutResponse struct {
	Status string `json:"status"`
}

// NewTemplatePut registers a sample PUT endpoint.
func NewTemplatePut(s *httpserver.Server) {
	fuegofw.Put(s.Manager, "/hello",
		func(c fuegofw.ContextWithBody[TemplatePutRequest]) (TemplatePutResponse, error) {
			body, err := c.Body()
			if err != nil {
				return TemplatePutResponse{}, err
			}
			return TemplatePutResponse{
				Status: body.Message + " updated",
			}, nil
		})
}
