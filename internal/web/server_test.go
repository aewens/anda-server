package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aewens/anda/pkg/reading"
)

var testServer *Server

func testCallback(server *Server) *Response {
	return &Response{
		Error: false,
		Name:  "test",
		Data:  "test",
	}
}

func TestCreate(t *testing.T) {
	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		t.Fatal("Missing 'CONFIG_PATH' environment variable")
	}

	defer Cleanup(t)
	config, err := reading.Config(configPath)
	if err != nil {
		t.Fatal("Could not read config file", err)
	}

	server, err := Create(config)
	if err != nil {
		t.Fatal(err)
	}

	testServer = server

	if server.Config != config {
		t.Error("Mutated value of config")
	}

	if server.Config != config {
		t.Error("Mutated value of config")
	}

	if _, ok := server.Routes["GET"]; !ok {
		t.Error("Routes missing 'GET'")
	}

	if _, ok := server.Routes["POST"]; !ok {
		t.Error("Routes missing 'POST'")
	}

	if _, ok := server.Routes["PUT"]; !ok {
		t.Error("Routes missing 'PUT'")
	}

	if _, ok := server.Routes["DELETE"]; !ok {
		t.Error("Routes missing 'DELETE'")
	}
}

func TestAddRoute(t *testing.T) {
	if testServer == nil {
		t.Fatal("Test server not defined")
	}

	testMethod := "GET"
	testRoute := "/api/test"

	testServer.AddRoute(testMethod, testRoute, testCallback)
	if _, ok := testServer.Routes[testMethod]; !ok {
		t.Fatalf("Missing method '%s'", testMethod)
	}

	if _, ok := testServer.Routes[testMethod][testRoute]; !ok {
		t.Fatalf("Missing route '%s'", testRoute)
	}
}

func TestWelcome(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	welcomeMethod := "GET"
	welcomeRoute := "/api"
	testServer.AddRoute(welcomeMethod, welcomeRoute, Welcome)
	handle := testServer.Routes[welcomeMethod][welcomeRoute]

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("/api status code '%v' instead of '%v'", status, http.StatusOK)
	}

	testError := false
	testName := "welcome"
	testData := "Hello, world!"
	template := "{\"error\":%t,\"name\":\"%s\",\"data\":\"%s\"}\n"
	testWelcome := fmt.Sprintf(template, testError, testName, testData)

	rrWelcome := rr.Body

	raw := make(map[string]interface{})
	err = json.NewDecoder(rrWelcome).Decode(&raw)
	if err != nil {
		t.Fatal(err)
	}

	rawError, ok := raw["error"]
	if !ok {
		t.Fatalf("Missing 'error'")
	}

	rawName, ok := raw["name"]
	if !ok {
		t.Fatalf("Missing 'name'")
	}

	rawData, ok := raw["data"]
	if !ok {
		t.Fatalf("Missing 'data'")
	}

	rrResponse := &Response{
		Error: rawError.(bool),
		Name:  rawName.(string),
		Data:  rawData.(string),
	}

	rrError := rrResponse.Error
	rrName := rrResponse.Name
	rrData := rrResponse.Data

	if rrError != testError {
		t.Errorf("For key 'Error' got '%v' instead of '%v'", rrError, testError)
	}

	if rrName != testName {
		t.Errorf("For key 'Name' got '%v' instead of '%v'", rrName, testName)
	}

	if rrData != testData {
		t.Errorf("For key 'Data' got '%v' instead of '%v'", rrData, testData)
	}

	if rrWelcome.String() != testWelcome {
		t.Errorf("/api returned '%v' instead of '%v'", rrWelcome, testWelcome)
	}
}
