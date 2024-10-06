package kitchen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

/*
KoksmatAutoPilotHostsTest

This file contains unit tests for the KoksmatAutoPilotHosts struct and its methods.
It tests the creation of new instances, loading from and saving to files, adding hosts,
setting default hosts, and retrieving host information.
*/

func TestNewKoksmatAutoPilotHosts(t *testing.T) {
	hosts := NewKoksmatAutoPilotHosts()
	if hosts == nil {
		t.Fatal("NewKoksmatAutoPilotHosts returned nil")
	}
	if hosts.config.Hosts == nil {
		t.Fatal("Hosts map is nil")
	}
}

func TestAddAndGetHost(t *testing.T) {
	hosts := NewKoksmatAutoPilotHosts()
	hosts.AddHost("test-host", "http://test.com", "test-key")

	host, err := hosts.GetHost("test-host")
	if err != nil {
		t.Fatalf("Failed to get host: %v", err)
	}
	if host.Href != "http://test.com" || host.Key != "test-key" {
		t.Fatalf("Host data doesn't match: got %+v", host)
	}

	_, err = hosts.GetHost("non-existent")
	if err == nil {
		t.Fatal("Expected error for non-existent host, got nil")
	}
}

func TestSetAndGetDefaultHost(t *testing.T) {
	hosts := NewKoksmatAutoPilotHosts()
	hosts.AddHost("default-host", "http://default.com", "default-key")

	err := hosts.SetDefaultHost("default-host")
	if err != nil {
		t.Fatalf("Failed to set default host: %v", err)
	}

	defaultHost, err := hosts.GetDefaultHost()
	if err != nil {
		t.Fatalf("Failed to get default host: %v", err)
	}
	if defaultHost.Href != "http://default.com" || defaultHost.Key != "default-key" {
		t.Fatalf("Default host data doesn't match: got %+v", defaultHost)
	}

	err = hosts.SetDefaultHost("non-existent")
	if err == nil {
		t.Fatal("Expected error when setting non-existent host as default, got nil")
	}
}

func TestSaveAndLoadFromFile(t *testing.T) {
	hosts := NewKoksmatAutoPilotHosts()
	hosts.AddHost("test-host", "http://test.com", "test-key")
	hosts.SetDefaultHost("test-host")

	tempDir, err := ioutil.TempDir("", "koksmat-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filename := filepath.Join(tempDir, "test-config.yaml")

	err = hosts.SaveToFile(filename)
	if err != nil {
		t.Fatalf("Failed to save to file: %v", err)
	}

	newHosts := NewKoksmatAutoPilotHosts()
	err = newHosts.LoadFromFile(filename)
	if err != nil {
		t.Fatalf("Failed to load from file: %v", err)
	}

	host, err := newHosts.GetHost("test-host")
	if err != nil {
		t.Fatalf("Failed to get host after loading: %v", err)
	}
	if host.Href != "http://test.com" || host.Key != "test-key" {
		t.Fatalf("Loaded host data doesn't match: got %+v", host)
	}

	defaultHost, err := newHosts.GetDefaultHost()
	if err != nil {
		t.Fatalf("Failed to get default host after loading: %v", err)
	}
	if defaultHost.Href != "http://test.com" || defaultHost.Key != "test-key" {
		t.Fatalf("Loaded default host data doesn't match: got %+v", defaultHost)
	}
}

func TestSaveToFileCreatePath(t *testing.T) {
	hosts := NewKoksmatAutoPilotHosts()
	hosts.AddHost("test-host", "http://test.com", "test-key")

	tempDir, err := ioutil.TempDir("", "koksmat-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	nonExistentDir := filepath.Join(tempDir, "non", "existent", "dir")
	filename := filepath.Join(nonExistentDir, "test-config.yaml")

	err = hosts.SaveToFile(filename)
	if err != nil {
		t.Fatalf("Failed to save to file in non-existent directory: %v", err)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("File was not created in the non-existent directory")
	}
}

func TestNiels(t *testing.T) {
	hosts := NewKoksmatAutoPilotHosts()
	auto, err := hosts.GetHost("auto")
	if err != nil {
		hosts.AddHost("auto", "http://test.com", "test-key")
	}
	auto, err = hosts.GetHost("auto")
	if err != nil {
		t.Fatalf("Failed to get host: %v", err)
	}
	if auto.Href != "http://test.com" || auto.Key != "test-key" {
		t.Fatalf("Host data doesn't match: got %+v", auto)
	}
	err = hosts.SaveToFile(".koksmat.autopilot.yaml")
	if err != nil {
		t.Fatalf("Failed to save to file in non-existent directory: %v", err)
	}
}
