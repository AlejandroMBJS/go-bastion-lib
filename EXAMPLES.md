# BastionRouter Examples Catalog

This catalog provides a comprehensive overview of the curated examples for the `goBastion Lib` framework. Each example is designed to be small, focused, and realistic, showcasing specific features, integrations, or best practices. They are copy-paste friendly, with comments explaining key choices and functionalities.

## Hello World + 404 Handler

**Folder:** `examples/hello-404/`  
**Focus:** Basic routing, custom JSON 404 handler, minimal setup.

You’ll learn how to:

- Bootstrap a `goBastion Lib` app.
- Register routes with `GET`.
- Plug in a custom 404 handler for JSON responses.

Run:

```bash
cd examples/hello-404
go run main.go
```

Test:

```bash
curl http://localhost:9876/
curl http://localhost:9876/does-not-exist
```

## JWT Auth (Core)

**Folder:** `examples/auth-jwt/`  
**Focus:** JSON Web Token (JWT) authentication, protected routes, token generation.

You’ll learn how to:

- Implement a login endpoint to issue JWTs.
- Use the `JWTAuth` middleware to protect routes.
- Extract user claims from the `router.Context`.

Run:

```bash
cd examples/auth-jwt
go run main.go
```

Test:

```bash
# Login to get a token
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"username":"testuser","password":"testpass"}' http://localhost:9876/api/auth/login | jq -r .access_token)

# Access protected route with token
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:9876/api/protected
```

## Auth + JWT + Docker + Postgres

**Folder:** `examples/auth-jwt-docker-postgres/`  
**Focus:** Full-stack authentication with JWT, PostgreSQL, and Docker Compose.

You’ll learn how to:

- Set up a `goBastion Lib` app with a PostgreSQL database using Docker Compose.
- Implement user registration, login, and JWT-based authentication.
- Perform basic database operations for user management.

Run:

```bash
cd examples/auth-jwt-docker-postgres
docker-compose up --build -d
go run main.go
```

Test:

```bash
# Register a user
curl -s -X POST -H "Content-Type: application/json" -d '{"username":"newuser","password":"newpass"}' http://localhost:9876/api/auth/register

# Login to get a token
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"username":"newuser","password":"newpass"}' http://localhost:9876/api/auth/login | jq -r .access_token)

# Access protected route
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:9876/api/profile
```

## CSRF Protection

**Folder:** `examples/csrf/`  
**Focus:** Cross-Site Request Forgery (CSRF) protection using double-submit cookie pattern.

You’ll learn how to:

- Integrate `CSRFMiddleware` into your application.
- Understand how CSRF tokens are set and validated.
- Simulate a client-side interaction with CSRF protection.

Run:

```bash
cd examples/csrf
go run main.go
```

Test (requires a client that handles cookies and headers, e.g., a browser or a more complex curl script):

```bash
# 1. Make a GET request to get the csrf_token cookie
#    (Simulate browser visiting the page)
#    This will set a 'csrf_token' cookie in the response.
#    You'd need to extract this cookie and send it back as X-CSRF-Token header.
#    For a simple curl test, you might manually set the cookie and header.
#    Example:
#    curl -c cookies.txt http://localhost:9876/
#    CSRF_TOKEN=$(grep csrf_token cookies.txt | awk '{print $7}')
#    curl -b cookies.txt -H "X-CSRF-Token: $CSRF_TOKEN" -X POST http://localhost:9876/api/submit
```

## GORM + PostgreSQL (CRUD API)

**Folder:** `examples/gorm-postgres-crud/`  
**Focus:** Full CRUD API with GORM ORM and PostgreSQL database.

You’ll learn how to:

- Integrate GORM with `goBastion Lib`.
- Define models and perform CRUD operations (Create, Read, Update, Delete).
- Use Docker Compose to manage a PostgreSQL database.

Run:

```bash
cd examples/gorm-postgres-crud
docker-compose up --build -d
go run main.go
```

Test:

```bash
# Create a product
curl -s -X POST -H "Content-Type: application/json" -d '{"name":"Laptop","price":1200.00}' http://localhost:9876/products

# Get all products
curl -s http://localhost:9876/products

# Get product by ID (assuming ID 1 was created)
curl -s http://localhost:9876/products/1

# Update a product
curl -s -X PUT -H "Content-Type: application/json" -d '{"name":"Gaming Laptop","price":1500.00}' http://localhost:9876/products/1

# Delete a product
curl -s -X DELETE http://localhost:9876/products/1
```

## File Upload

**Folder:** `examples/file-upload/`  
**Focus:** Handling multipart form data for file uploads and serving static files.

You’ll learn how to:

- Process `multipart/form-data` requests.
- Save uploaded files securely to disk.
- Serve static files from a designated directory.

Run:

```bash
cd examples/file-upload
go run main.go
```

Test:

```bash
# Upload a file (replace 'your_file.txt' with an actual file)
curl -s -X POST -F "file=@your_file.txt" http://localhost:9876/upload

# Access the uploaded file (check server output for filename)
# Example: curl http://localhost:9876/files/uploaded_file_name.txt
```

## WebSocket Echo Server

**Folder:** `examples/websocket/`  
**Focus:** Real-time communication using WebSockets.

You’ll learn how to:

- Upgrade an HTTP connection to a WebSocket.
- Send and receive messages over a WebSocket connection.
- Implement a simple echo server.

Run:

```bash
cd examples/websocket
go run main.go
```

Test (requires a WebSocket client, e.g., `wscat` or a browser's developer console):

```bash
# Using wscat (install with: npm install -g wscat)
wscat -c ws://localhost:9876/ws
> Hello
< Hello
```

## OpenAPI + Swagger Docs

**Folder:** `examples/openapi-swagger/`  
**Focus:** Generating and serving OpenAPI (Swagger) documentation.

You’ll learn how to:

- Serve a static `openapi.yaml` file.
- Integrate Swagger UI for interactive API documentation.
- Document your API endpoints using OpenAPI specifications.

Run:

```bash
cd examples/openapi-swagger
go run main.go
```

Test:

```bash
# View the raw OpenAPI spec
curl http://localhost:9876/openapi.yaml

# Open in browser to see Swagger UI
# http://localhost:9876/swagger/
```

## Testing (Unit & Integration)

**Folder:** `examples/testing/`  
**Focus:** Best practices for unit and integration testing `goBastion Lib` applications.

You’ll learn how to:

- Write unit tests for individual handler functions.
- Use `net/http/httptest` for in-memory integration tests of your routes and middleware.
- Structure your tests for clarity and maintainability.

Run:

```bash
cd examples/testing
go test ./... -v
```

## Docker + Nginx Reverse Proxy

**Folder:** `examples/docker-nginx/`  
**Focus:** Deploying a `goBastion Lib` application behind an Nginx reverse proxy using Docker Compose.

You’ll learn how to:

- Containerize your Go application with Docker.
- Configure Nginx as a reverse proxy to forward requests to your Go app.
- Set up a production-like deployment environment with Docker Compose.

Run:

```bash
cd examples/docker-nginx
docker-compose up --build -d
```

Test:

```bash
# Access the Go app via Nginx
curl http://localhost:80/api/hello
```

## Graceful Shutdown

**Folder:** `examples/graceful-shutdown/`  
**Focus:** Implementing graceful server shutdown for resilient deployments.

You’ll learn how to:

- Handle OS signals (SIGINT, SIGTERM) to initiate a graceful shutdown.
- Allow active requests to complete before the server fully stops.
- Ensure your application cleans up resources properly.

Run:

```bash
cd examples/graceful-shutdown
go run main.go
# In another terminal, send a SIGINT (Ctrl+C) to the running process
# Observe the server logs for graceful shutdown messages
```

Test:

```bash
# Start the server in one terminal
# go run main.go

# In another terminal, send a request that takes some time
# curl http://localhost:9876/long-task

# While the long-task is running, press Ctrl+C in the server terminal.
# The server should log "Shutting down server gracefully..." and the long-task request should still complete.
```

---

👉 For a visual overview and design-focused tour, see `docs/index.html`.