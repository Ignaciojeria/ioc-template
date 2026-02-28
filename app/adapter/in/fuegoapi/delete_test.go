package fuegoapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"archetype/app/shared/configuration"
	"archetype/app/shared/infrastructure/httpserver"
)

func TestNewTemplateDelete(t *testing.T) {
	conf := configuration.Conf{
		PORT:         "8089",
		PROJECT_NAME: "test",
		VERSION:      "v1",
	}

	server, err := httpserver.NewServer(conf)
	if err != nil {
		t.Fatalf("unexpected error creating server: %v", err)
	}

	NewTemplateDelete(server)

	req, err := http.NewRequest(http.MethodDelete, "/hello/123", nil)
	if err != nil {
		t.Fatalf("unexpected error building request: %v", err)
	}

	recorder := httptest.NewRecorder()
	server.Manager.Mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusOK)
	}

	expectedBody := "123 deleted"
	if recorder.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %s want %s", recorder.Body.String(), expectedBody)
	}
}
