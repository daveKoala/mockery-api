// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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

func main() {
	configFile := flag.String("config", "config.json", "Path to configuration file")
	outputFile := flag.String("output", "ENDPOINTS.md", "Output file for endpoints documentation")
	flag.Parse()

	// Load config
	data, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// Open output file
	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()

	baseURL := fmt.Sprintf("http://localhost:%d", config.Server.Port)

	// Write markdown header
	fmt.Fprintln(f, "# API Endpoints")
	fmt.Fprintln(f, "")
	fmt.Fprintf(f, "> Auto-generated from `%s`\n", *configFile)
	fmt.Fprintln(f, "")
	fmt.Fprintf(f, "**Base URL:** `%s`\n", baseURL)
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "---")
	fmt.Fprintln(f, "")

	// Add health check
	fmt.Fprintln(f, "## Health Check")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "**GET** `/_health`")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "```bash")
	fmt.Fprintf(f, "curl %s/_health\n", baseURL)
	fmt.Fprintln(f, "```")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "---")
	fmt.Fprintln(f, "")

	// Generate documentation for each route
	for _, route := range config.Routes {
		writeRouteDoc(f, route, baseURL)
	}

	fmt.Printf("Generated documentation for %d endpoints in %s\n", len(config.Routes)+1, *outputFile)
}

func writeRouteDoc(f *os.File, route Route, baseURL string) {
	// Convert path parameters to examples
	examplePath := convertPathToExample(route.Path)

	// Header with method and path
	fmt.Fprintf(f, "## %s %s\n", route.Method, route.Path)
	fmt.Fprintln(f, "")

	// Auth requirements
	if route.RequiresAuth {
		fmt.Fprintf(f, "🔒 **Requires Authentication:** `%s` header\n", route.AuthHeader)
		fmt.Fprintln(f, "")
	}

	// Response info
	fmt.Fprintf(f, "**Response:** `%d`\n", route.Response.Status)
	fmt.Fprintln(f, "")

	// Build curl command
	var curlParts []string
	curlParts = append(curlParts, "curl")

	if route.Method != "GET" && route.Method != "HEAD" {
		curlParts = append(curlParts, fmt.Sprintf("-X %s", route.Method))
	}

	if route.RequiresAuth {
		authValue := "YOUR_TOKEN_HERE"
		if route.AuthHeader == "Authorization" {
			authValue = "Bearer " + authValue
		}
		curlParts = append(curlParts, fmt.Sprintf("-H \"%s: %s\"", route.AuthHeader, authValue))
	}

	if route.Method == "HEAD" {
		curlParts = append(curlParts, "-I")
	}

	curlParts = append(curlParts, fmt.Sprintf("%s%s", baseURL, examplePath))

	// Curl example
	fmt.Fprintln(f, "```bash")
	fmt.Fprintln(f, strings.Join(curlParts, " "))
	fmt.Fprintln(f, "```")
	fmt.Fprintln(f, "")

	// Example response body (if not empty/null)
	if route.Response.Body != nil && route.Response.Status != 204 {
		fmt.Fprintln(f, "<details>")
		fmt.Fprintln(f, "<summary>Example Response</summary>")
		fmt.Fprintln(f, "")
		fmt.Fprintln(f, "```json")
		bodyJSON, _ := json.MarshalIndent(route.Response.Body, "", "  ")
		fmt.Fprintln(f, string(bodyJSON))
		fmt.Fprintln(f, "```")
		fmt.Fprintln(f, "</details>")
		fmt.Fprintln(f, "")
	}

	fmt.Fprintln(f, "---")
	fmt.Fprintln(f, "")
}

func convertPathToExample(path string) string {
	// Replace {param} with example values
	replacements := map[string]string{
		"{id}":        "123",
		"{userId}":    "456",
		"{productId}": "789",
		"{orderId}":   "order-123",
		"{itemId}":    "item-456",
	}

	result := path
	for param, example := range replacements {
		result = strings.ReplaceAll(result, param, example)
	}

	// Replace any remaining {something} with generic value
	for strings.Contains(result, "{") {
		start := strings.Index(result, "{")
		end := strings.Index(result, "}")
		if end > start {
			result = result[:start] + "example-value" + result[end+1:]
		} else {
			break
		}
	}

	return result
}
