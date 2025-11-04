# Mockery API

A simple, pragmatic API mocker built in Go that serves mock responses based on a JSON configuration file.

## Features

- Simple JSON-based configuration
- Support for all common HTTP methods (GET, POST, PUT, DELETE, PATCH, HEAD)
- Basic auth header validation
- Static response mocking
- Clean stdout logging
- No external dependencies (pure Go stdlib)

## Getting Started

### Quick Start with Make

```bash
# Build and run
make run

# Or build and start in background
make start

# Stop background server
make stop

# Check server status
make status

# View logs (for background server)
make logs

# Run tests
make test

# Clean everything
make clean
```

### Manual Build and Run

```bash
# Build
go build -o mockery-api

# Run
./mockery-api

# Run with custom config
./mockery-api -config path/to/your/config.json
```

### Available Make Commands

- `make help` - Show all available commands
- `make build` - Build the binary
- `make run` - Build and run in foreground
- `make start` - Build and start in background (logs to server.log)
- `make stop` - Stop background server
- `make restart` - Restart background server
- `make status` - Check if server is running
- `make logs` - Tail server logs
- `make test` - Run basic API tests
- `make curls` - Generate ENDPOINTS.md from config file
- `make clean` - Remove binary, logs, and PID file

You can specify a custom config file:
```bash
make run CONFIG=my-config.json
make start CONFIG=my-config.json
make curls CONFIG=my-config.json
```

### Generate Endpoint Documentation

Automatically generate markdown documentation with curl examples from your config:

```bash
make curls
```

This creates `ENDPOINTS.md` with:
- Clean, readable documentation for all endpoints
- Path parameters converted to example values
- Auth requirements clearly marked
- Copy-paste ready curl commands
- Example response bodies (collapsible)

The markdown format makes it easy to:
- Read in any editor or on GitHub
- Convert to other formats (Postman, HTTP files, etc.)
- Share with your team
- Commit as API documentation

## Configuration Format

The configuration file uses a simple JSON structure:

```json
{
  "server": {
    "port": 3000
  },
  "routes": [
    {
      "path": "/api/users",
      "method": "GET",
      "requiresAuth": true,
      "authHeader": "Authorization",
      "response": {
        "status": 200,
        "headers": {
          "X-Custom-Header": "value"
        },
        "body": {
          "users": []
        }
      }
    }
  ]
}
```

### Configuration Fields

#### Server
- `port` (required): Port number to run the server on

#### Route
- `path` (required): Path to match. Supports path parameters using `{paramName}` syntax
  - Static: `/api/users`
  - With parameters: `/api/products/{id}` or `/api/orders/{orderId}/items/{itemId}`
- `method` (required): HTTP method (GET, POST, PUT, DELETE, PATCH, HEAD)
- `requiresAuth` (optional): Whether to check for auth header (default: false)
- `authHeader` (optional): Name of the auth header to check (required if `requiresAuth` is true)
- `response` (required): Response configuration

#### Response
- `status` (required): HTTP status code to return
- `headers` (optional): Custom response headers
- `body` (optional): JSON response body (can be null for 204 responses)

### Path Parameters

Path parameters allow you to define a single route that matches multiple URLs. Use `{paramName}` syntax:

```json
{
  "path": "/api/products/{id}",
  "method": "GET",
  "response": {
    "status": 200,
    "body": {
      "id": 101,
      "name": "Product"
    }
  }
}
```

This route will match:
- `/api/products/101`
- `/api/products/abc`
- `/api/products/xyz123`

You can use multiple parameters:
```json
{
  "path": "/api/orders/{orderId}/items/{itemId}",
  "method": "GET"
}
```

**Note:** Path parameter values are not currently extracted or used in responses. The same static response is returned regardless of the parameter value. This is perfect for development where you just need to avoid hitting expensive APIs.

## Examples

### Testing with curl

```bash
# Health check
curl http://localhost:3000/_health

# GET without auth
curl http://localhost:3000/api/products

# GET with auth (will fail without header)
curl http://localhost:3000/api/users

# GET with auth header
curl -H "Authorization: Bearer token123" http://localhost:3000/api/users

# POST with auth
curl -X POST -H "Authorization: Bearer token123" http://localhost:3000/api/users

# GET with path parameter (matches /api/products/{id})
curl http://localhost:3000/api/products/101
curl http://localhost:3000/api/products/202
curl http://localhost:3000/api/products/anything

# GET user by ID with auth (matches /api/users/{userId})
curl -H "Authorization: Bearer token" http://localhost:3000/api/users/123

# PUT with path parameter and auth
curl -X PUT -H "Authorization: Bearer token123" http://localhost:3000/api/products/999

# Multiple path parameters (matches /api/orders/{orderId}/items/{itemId})
curl -H "Authorization: Bearer token" http://localhost:3000/api/orders/12345/items/67890

# PATCH with custom auth header
curl -X PATCH -H "X-API-Key: mykey" http://localhost:3000/api/settings

# HEAD request
curl -I -X HEAD http://localhost:3000/api/status
```

## Logging

The server logs all incoming requests to stdout with the following information:
- HTTP method and path
- Route matching status
- Auth validation status (if required)
- Response status code

Example log output:
```
[GET] /api/users
  ✓ Matched route: GET /api/users
  ✓ Auth header 'Authorization' present
  ✓ Response sent: 200
```

## Built-in Endpoints

- `/_health` - Health check endpoint that returns `{"status":"ok","message":"mockery-api is running"}`

## Notes

- Route matching is exact (no wildcard or regex support)
- Auth validation only checks if the header exists, not its value
- Responses are always returned as JSON with `Content-Type: application/json`
- The server must be restarted to pick up config changes
