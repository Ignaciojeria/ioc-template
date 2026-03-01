# fuegoapi-controllers

> REST API controllers scaffolded with Fuego (GET, POST, PUT, PATCH, DELETE)

## app/adapter/in/fuegoapi/get.go

```go
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
```

---

## app/adapter/in/fuegoapi/post.go

```go
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
```

---

## app/adapter/in/fuegoapi/put.go

```go
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
```

---

## app/adapter/in/fuegoapi/patch.go

```go
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
```

---

## app/adapter/in/fuegoapi/delete.go

```go
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
```
