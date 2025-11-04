package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the main configuration structure
type Config struct {
	Server ServerConfig `json:"server"`
	Routes []Route      `json:"routes"`
}

// ServerConfig holds server-specific settings
type ServerConfig struct {
	Port int `json:"port"`
}

// Route represents a single API endpoint configuration
type Route struct {
	Path         string            `json:"path"`
	Method       string            `json:"method"`
	RequiresAuth bool              `json:"requiresAuth"`
	AuthHeader   string            `json:"authHeader"`
	Response     Response          `json:"response"`
}

// Response represents the mock response configuration
type Response struct {
	Status  int                    `json:"status"`
	Headers map[string]string      `json:"headers,omitempty"`
	Body    interface{}            `json:"body"`
}

// LoadConfig reads and parses the configuration file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// validateConfig performs basic validation on the configuration
func validateConfig(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", config.Server.Port)
	}

	validMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
		"PATCH":  true,
		"HEAD":   true,
	}

	for i, route := range config.Routes {
		if route.Path == "" {
			return fmt.Errorf("route %d: path cannot be empty", i)
		}
		if !validMethods[route.Method] {
			return fmt.Errorf("route %d: invalid method %s", i, route.Method)
		}
		if route.RequiresAuth && route.AuthHeader == "" {
			return fmt.Errorf("route %d: authHeader required when requiresAuth is true", i)
		}
	}

	return nil
}
