package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Load configuration
	log.Printf("Loading configuration from: %s", *configFile)
	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Configuration loaded successfully")
	log.Printf("  - Port: %d", config.Server.Port)
	log.Printf("  - Routes: %d configured", len(config.Routes))

	// Create handler with configured routes
	handler := NewMockHandler(config.Routes)

	// Setup HTTP server with mux
	mux := http.NewServeMux()

	// Add health check endpoint
	mux.HandleFunc("/_health", healthCheckHandler)

	// Add catch-all handler for mock routes
	mux.Handle("/", handler)

	// Server address
	addr := fmt.Sprintf(":%d", config.Server.Port)

	// Start server
	log.Printf("Starting mockery-api server on http://localhost%s", addr)
	log.Printf("Health check available at: http://localhost%s/_health", addr)
	log.Println("Press Ctrl+C to stop")
	log.Println("---")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}
