---
name: fuego-controllers
description: Boilerplate for defining REST APIs using Fuego
---

## Overview

All incoming REST interfaces MUST be implemented using the `github.com/go-fuego/fuego` framework.
You MUST NOT use `net/http` handlers natively or any other framework like Gin, Fiber, or Echo.

## Structure of a Fuego Controller

A Fuego controller lives in `app/adapter/in/fuegoapi/`. It uses the IoC container to receive dependencies (like the Fuego `*httpserver.Server` and a UseCase).

> [!NOTE]
> Read the file `app/adapter/in/fuegoapi/controller.go` (if it exists in your workspace) to see a living implementation of this boilerplate.

### Mandatory Template

Every controller you write MUST follow this exact structure:

```go
package fuegoapi

import (
	"context"
	"net/http"

	"archetype/app/shared/infrastructure/httpserver"

	"github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
)

// 1. Register the Controller with IoC
var _ = ioc.Register(NewMyCustomController)

type MyCustomController struct {
	// Add your use cases or dependencies here
}

// 2. Inject the httpserver.Server instance
func NewMyCustomController(server *httpserver.Server /*, usecase MyUseCase */) (*MyCustomController, error) {
	c := &MyCustomController{}
	
	// 3. Register the route directly onto the Fuego server instance
	fuego.Post(server.Manager, "/api/v1/my-resource", c.HandleRequest)
	
	return c, nil
}

// 4. Define Fuego Schema types for validation
type MyRequest struct {
	Name string `json:"name" validate:"required"`
}

type MyResponse struct {
	Message string `json:"message"`
}

// 5. Build the Handler conforming to Fuego's signature
func (c *MyCustomController) HandleRequest(ctx fuego.ContextWithBody[MyRequest]) (MyResponse, error) {
	body, err := ctx.Body()
	if err != nil {
		return MyResponse{}, fuego.BadRequestError{Err: err}
	}
	
	// Delegate to your core/usecase here
	
	return MyResponse{Message: "Hello " + body.Name}, nil
}
```

## Validation & Errors

- **Validation**: Use `validate:"required"` tags on your DTOs. Fuego catches them automatically.
- **Errors**: Return Fuego standardized errors (e.g., `fuego.BadRequestError{Err: err}`) instead of manually writing HTTP status codes.
- **Context**: Use `ctx.Context()` when passing context down to databases or event buses, NOT the `fuego.ContextWithBody` object itself.
