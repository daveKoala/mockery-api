package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// MockHandler handles incoming HTTP requests and matches them against configured routes
type MockHandler struct {
	routes []Route
}

// NewMockHandler creates a new handler with the given routes
func NewMockHandler(routes []Route) *MockHandler {
	return &MockHandler{
		routes: routes,
	}
}

// ServeHTTP implements the http.Handler interface
func (h *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log incoming request
	log.Printf("[%s] %s", r.Method, r.URL.Path)

	// Find matching route
	route := h.findRoute(r.Method, r.URL.Path)
	if route == nil {
		log.Printf("  ✗ No route matched")
		http.NotFound(w, r)
		return
	}

	log.Printf("  ✓ Matched route: %s %s", route.Method, route.Path)

	// Check auth if required
	if route.RequiresAuth {
		authValue := r.Header.Get(route.AuthHeader)
		if authValue == "" {
			log.Printf("  ✗ Auth failed: missing header '%s'", route.AuthHeader)
			http.Error(w, "Unauthorized: missing auth header", http.StatusUnauthorized)
			return
		}
		log.Printf("  ✓ Auth header '%s' present", route.AuthHeader)
	}

	// Set custom response headers if configured
	if route.Response.Headers != nil {
		for key, value := range route.Response.Headers {
			w.Header().Set(key, value)
		}
	}

	// Always set Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write status code
	w.WriteHeader(route.Response.Status)

	// Write response body
	if route.Response.Body != nil {
		if err := json.NewEncoder(w).Encode(route.Response.Body); err != nil {
			log.Printf("  ✗ Error encoding response: %v", err)
			return
		}
	}

	log.Printf("  ✓ Response sent: %d", route.Response.Status)
}

// findRoute searches for a matching route based on method and path
func (h *MockHandler) findRoute(method, path string) *Route {
	for i := range h.routes {
		route := &h.routes[i]
		if route.Method == method && route.Path == path {
			return route
		}
	}
	return nil
}

// healthCheckHandler provides a simple health check endpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","message":"mockery-api is running"}`)
}
