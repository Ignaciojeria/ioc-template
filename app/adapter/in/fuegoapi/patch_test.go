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

func TestNewTemplatePatch(t *testing.T) {
	conf := configuration.Conf{
		PORT:         "8087",
		PROJECT_NAME: "test",
		VERSION:      "v1",
	}

	server, err := httpserver.NewServer(conf, nil)
	if err != nil {
		t.Fatalf("unexpected error creating server: %v", err)
	}

	NewTemplatePatch(server)

	reqBody, _ := json.Marshal(TemplatePatchRequest{Message: "Einar"})
	req, err := http.NewRequest(http.MethodPatch, "/hello", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("unexpected error building request: %v", err)
	}

	recorder := httptest.NewRecorder()
	server.Manager.Mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusOK)
	}

	expectedBody := `{"status":"Einar patched"}` + "\n"
	if recorder.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %s want %s", recorder.Body.String(), expectedBody)
	}
}

func TestNewTemplatePatch_InvalidBody(t *testing.T) {
	conf := configuration.Conf{
		PORT:         "8088",
		PROJECT_NAME: "test-err",
		VERSION:      "v1",
	}

	server, err := httpserver.NewServer(conf, nil)
	if err != nil {
		t.Fatalf("unexpected error creating server: %v", err)
	}

	NewTemplatePatch(server)

	req, err := http.NewRequest(http.MethodPatch, "/hello", bytes.NewBuffer([]byte("{invalid json")))
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
