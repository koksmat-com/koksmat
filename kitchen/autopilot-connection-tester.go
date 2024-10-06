package kitchen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/*
AutopilotConnectionTester

This file contains the implementation of AutopilotConnectionTester, which
provides functionality to test connections with Koksmat AutoPilot hosts.
It supports registering connections, pinging, posting requests, and getting status.
*/

// AutopilotConnectionTester represents a connection tester for Koksmat AutoPilot hosts
type AutopilotConnectionTester struct {
	config     KoksmatAutoPilotHostConfig
	jwtToken   string
	httpClient *http.Client
}

// NewAutopilotConnectionTester creates a new AutopilotConnectionTester instance
func NewAutopilotConnectionTester(config KoksmatAutoPilotHostConfig, jwtToken string) *AutopilotConnectionTester {
	return &AutopilotConnectionTester{
		config:   config,
		jwtToken: jwtToken,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// RegisterConnection registers a new connection with the AutoPilot host
func (a *AutopilotConnectionTester) RegisterConnection(key, clientSecret string) error {
	url := fmt.Sprintf("%s/api/autopilot/register", a.config.Href)

	// Create the request body
	requestBody := struct {
		Key          string `json:"key"`
		ClientSecret string `json:"clientSecret"`
	}{
		Key:          key,
		ClientSecret: clientSecret,
	}

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Send the POST request with the JSON body
	return a.sendRequest("POST", url, jsonBody)
}

// Ping sends a ping request to the AutoPilot host
func (a *AutopilotConnectionTester) Ping() error {
	url := fmt.Sprintf("%s/api/autopilot/ping", a.config.Href)
	return a.sendRequest("GET", url, nil)
}

// PostRequest sends a POST request to the AutoPilot host
func (a *AutopilotConnectionTester) PostRequest(body []byte) error {
	url := fmt.Sprintf("%s/api/autopilot/request", a.config.Href)
	return a.sendRequest("POST", url, body)
}

// GetStatus retrieves the status from the AutoPilot host
func (a *AutopilotConnectionTester) GetStatus() (string, error) {
	url := fmt.Sprintf("%s/api/autopilot/status", a.config.Href)
	resp, err := a.sendRequestWithResponse("GET", url, nil)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// sendRequest sends an HTTP request to the specified URL
func (a *AutopilotConnectionTester) sendRequest(method, url string, body []byte) error {
	_, err := a.sendRequestWithResponse(method, url, body)
	return err
}

// sendRequestWithResponse sends an HTTP request and returns the response body
func (a *AutopilotConnectionTester) sendRequestWithResponse(method, url string, body []byte) ([]byte, error) {

	var err error

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.jwtToken))

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
