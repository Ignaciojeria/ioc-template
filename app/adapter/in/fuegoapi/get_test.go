package fuegoapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"archetype/app/shared/configuration"
	"archetype/app/shared/infrastructure/httpserver"
)

func TestNewTemplateGet(t *testing.T) {
	conf := configuration.Conf{
		PORT:         "8083",
		PROJECT_NAME: "test",
		VERSION:      "v1",
	}

	server, err := httpserver.NewServer(conf, nil)
	if err != nil {
		t.Fatalf("unexpected error creating server: %v", err)
	}

	// Register the endpoint using our component logic
	NewTemplateGet(server)

	// Create a simulated test request targeting the endpoint
	req, err := http.NewRequest(http.MethodGet, "/hello", nil)
	if err != nil {
		t.Fatalf("unexpected error building request: %v", err)
	}

	// Record the response directly from the Fuego router's underlying Mux
	recorder := httptest.NewRecorder()
	server.Manager.Mux.ServeHTTP(recorder, req)

	// Assert the status code
	if recorder.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusOK)
	}

	// Fuego string handling directly outputs the string content
	expectedBody := "Hello from Einar IoC Fuego Template!"
	if recorder.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %s want %s", recorder.Body.String(), expectedBody)
	}
}
