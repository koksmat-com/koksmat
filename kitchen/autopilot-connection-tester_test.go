package kitchen

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
AutopilotConnectionTesterTest

This file contains unit tests for the AutopilotConnectionTester struct and its methods.
It tests the creation of new instances, registering connections, pinging,
posting requests, and getting status using a mock HTTP server.
*/

func TestNewAutopilotConnectionTester(t *testing.T) {
	config := KoksmatAutoPilotHostConfig{
		Href: "http://example.com",
		Key:  "test-key",
	}
	jwtToken := "test-token"

	tester := NewAutopilotConnectionTester(config, jwtToken)
	if tester == nil {
		t.Fatal("NewAutopilotConnectionTester returned nil")
	}
	if tester.config != config {
		t.Fatalf("Config mismatch. Expected %v, got %v", config, tester.config)
	}
	if tester.jwtToken != jwtToken {
		t.Fatalf("JWT token mismatch. Expected %s, got %s", jwtToken, tester.jwtToken)
	}
}

func TestRegisterConnection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/autopilot/register" {
			t.Errorf("Expected to request '/api/autopilot/register', got: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := KoksmatAutoPilotHostConfig{
		Href: server.URL,
		Key:  "test-key",
	}
	tester := NewAutopilotConnectionTester(config, "test-token")

	err := tester.RegisterConnection()
	if err != nil {
		t.Fatalf("RegisterConnection failed: %v", err)
	}
}

func TestPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/autopilot/ping" {
			t.Errorf("Expected to request '/api/autopilot/ping', got: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := KoksmatAutoPilotHostConfig{
		Href: server.URL,
		Key:  "test-key",
	}
	tester := NewAutopilotConnectionTester(config, "test-token")

	err := tester.Ping()
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}

func TestPostRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/autopilot/request" {
			t.Errorf("Expected to request '/api/autopilot/request', got: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got: %s", r.Method)
		}
		var data map[string]interface{}
		json.NewDecoder(r.Body).Decode(&data)
		if data["key"] != "value" {
			t.Errorf("Expected request body to contain {'key': 'value'}, got: %v", data)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := KoksmatAutoPilotHostConfig{
		Href: server.URL,
		Key:  "test-key",
	}
	tester := NewAutopilotConnectionTester(config, "test-token")

	err := tester.PostRequest(map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("PostRequest failed: %v", err)
	}
}

func TestGetStatus2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/autopilot/status" {
			t.Errorf("Expected to request '/api/autopilot/status', got: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	config := KoksmatAutoPilotHostConfig{
		Href: server.URL,
		Key:  "test-key",
	}
	tester := NewAutopilotConnectionTester(config, "test-token")

	status, err := tester.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status != "OK" {
		t.Fatalf("Expected status 'OK', got: %s", status)
	}
}

func TestAuthorizationHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expectedHeader := "Bearer test-token"
		if authHeader != expectedHeader {
			t.Errorf("Expected Authorization header '%s', got: '%s'", expectedHeader, authHeader)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := KoksmatAutoPilotHostConfig{
		Href: server.URL,
		Key:  "test-key",
	}
	tester := NewAutopilotConnectionTester(config, "test-token")

	err := tester.Ping()
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}
