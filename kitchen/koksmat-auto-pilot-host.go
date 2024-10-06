package kitchen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

/*
KoksmatAutoPilotHosts

This file contains the implementation of KoksmatAutoPilotHosts, which manages
configurations for multiple Koksmat AutoPilot Hosts. It provides functionality
to load and save configurations from/to YAML files, add and retrieve host
information, and manage the default host.
*/

// KoksmatAutoPilotHostConfig represents the configuration for a single KoksmatAutoPilotHost
type KoksmatAutoPilotHostConfig struct {
	Href string `yaml:"href"`
	Key  string `yaml:"key"`
}

// KoksmatAutoPilotHostsConfig represents the entire configuration
type KoksmatAutoPilotHostsConfig struct {
	Hosts       map[string]KoksmatAutoPilotHostConfig `yaml:"hosts"`
	DefaultHost string                                `yaml:"default_host"`
}

// KoksmatAutoPilotHosts is the main struct to handle KoksmatAutoPilotHost configurations
type KoksmatAutoPilotHosts struct {
	config KoksmatAutoPilotHostsConfig
}

// NewKoksmatAutoPilotHosts creates a new KoksmatAutoPilotHosts instance
func NewKoksmatAutoPilotHosts() *KoksmatAutoPilotHosts {
	return &KoksmatAutoPilotHosts{
		config: KoksmatAutoPilotHostsConfig{
			Hosts: make(map[string]KoksmatAutoPilotHostConfig),
		},
	}
}

// LoadFromFile loads the configuration from a YAML file
func (k *KoksmatAutoPilotHosts) LoadFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &k.config)
	if err != nil {
		return err
	}

	return nil
}

// GetHost returns the configuration for a specific KoksmatAutoPilotHost
func (k *KoksmatAutoPilotHosts) GetHost(tag string) (KoksmatAutoPilotHostConfig, error) {
	if host, ok := k.config.Hosts[tag]; ok {
		return host, nil
	}
	return KoksmatAutoPilotHostConfig{}, fmt.Errorf("KoksmatAutoPilotHost with tag '%s' not found", tag)
}

// GetDefaultHost returns the configuration for the default KoksmatAutoPilotHost
func (k *KoksmatAutoPilotHosts) GetDefaultHost() (KoksmatAutoPilotHostConfig, error) {
	return k.GetHost(k.config.DefaultHost)
}

// AddHost adds a new KoksmatAutoPilotHost to the configuration
func (k *KoksmatAutoPilotHosts) AddHost(tag string, href string, key string) {
	k.config.Hosts[tag] = KoksmatAutoPilotHostConfig{
		Href: href,
		Key:  key,
	}
}

// SetDefaultHost sets the default KoksmatAutoPilotHost
func (k *KoksmatAutoPilotHosts) SetDefaultHost(tag string) error {
	if _, ok := k.config.Hosts[tag]; !ok {
		return fmt.Errorf("KoksmatAutoPilotHost with tag '%s' not found", tag)
	}
	k.config.DefaultHost = tag
	return nil
}

// SaveToFile saves the current configuration to a YAML file
// It creates the directory path if it doesn't exist
func (k *KoksmatAutoPilotHosts) SaveToFile(filename string) error {
	// Ensure the directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Marshal the config to YAML
	data, err := yaml.Marshal(&k.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write the file
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
