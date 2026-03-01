package fuegoapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"archetype/app/shared/configuration"
	"archetype/app/shared/infrastructure/httpserver"
)

func TestNewTemplatePut(t *testing.T) {
	conf := configuration.Conf{
		PORT:         "8085",
		PROJECT_NAME: "test",
		VERSION:      "v1",
	}

	server, err := httpserver.NewServer(conf, nil)
	if err != nil {
		t.Fatalf("unexpected error creating server: %v", err)
	}

	NewTemplatePut(server)

	reqBody, _ := json.Marshal(TemplatePutRequest{Message: "Einar"})
	req, err := http.NewRequest(http.MethodPut, "/hello", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("unexpected error building request: %v", err)
	}

	recorder := httptest.NewRecorder()
	server.Manager.Mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusOK)
	}

	expectedBody := `{"status":"Einar updated"}` + "\n"
	if recorder.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %s want %s", recorder.Body.String(), expectedBody)
	}
}

func TestNewTemplatePut_InvalidBody(t *testing.T) {
	conf := configuration.Conf{
		PORT:         "8086",
		PROJECT_NAME: "test-err",
		VERSION:      "v1",
	}

	server, err := httpserver.NewServer(conf, nil)
	if err != nil {
		t.Fatalf("unexpected error creating server: %v", err)
	}

	NewTemplatePut(server)

	req, err := http.NewRequest(http.MethodPut, "/hello", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("unexpected error building request: %v", err)
	}

	recorder := httptest.NewRecorder()
	server.Manager.Mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("handler should return bad request: got %v want %v", recorder.Code, http.StatusBadRequest)
	}
}
