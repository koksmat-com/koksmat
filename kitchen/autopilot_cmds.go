package kitchen

import (
	"errors"

	"github.com/koksmat-com/koksmat/autopilot"
)

func AutoPilotRun(id string) (string, error) {
	// Run the auto pilot mode

	autopilot.Run(id)
	return "", errors.New("Not implemented")
}

func HasAutopilotConnection() (bool, error) {
	// Check if the user has a connection to Koksmat Studio
	return false, errors.New("Not implemented")
}

func AddAutopilotConnection(token string, serverUrl string) (string, error) {
	// Make a connection to Koksmat Studio
	return "", errors.New("Not implemented")
}

func RemoveAutopilotConnection(id string) error {
	// Remove the connection to Koksmat Studio
	return errors.New("Not implemented")
}

func ListAutopilotConnections() ([]string, error) {
	// List all connections to Koksmat Studio
	return nil, errors.New("Not implemented")
}

func GetAutopilotConnection(id string) (string, error) {
	// Get a connection to Koksmat Studio
	return "", errors.New("Not implemented")
}

func UpdateAutopilotConnection(id string, token string, serverUrl string) (string, error) {
	// Update a connection to Koksmat Studio
	return "", errors.New("Not implemented")
}

func SetAutopilotDefaultConnection(id string) error {
	// Set the default connection to Koksmat Studio
	return errors.New("Not implemented")
}

func GetAutopilotDefaultConnection() (string, error) {
	// Get the default connection to Koksmat Studio
	return "auto", nil
	//return "", errors.New("Not implemented")
}
