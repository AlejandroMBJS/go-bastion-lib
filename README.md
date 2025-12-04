# goBastion Lib

**Secure, opinionated HTTP routing for Go JSON APIs.**

[![Go Reference](https://pkg.go.dev/badge/github.com/alejandrombjs/go-bastion-lib.svg)](https://pkg.go.dev/github.com/alejandrombjs/go-bastion-lib)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/alejandrombjs/go-bastion-lib)](https://goreportcard.com/report/github.com/alejandrombjs/go-bastion-lib)

`goBastion Lib` is a robust and opinionated Go library designed to streamline the development of secure, high-performance JSON APIs. It provides a foundational framework that handles common web concerns like routing, middleware, and security, allowing developers to focus on their core business logic.

## Why `goBastion Lib` instead of X?

*   **JSON-only Focus:** Explicitly designed for building HTTP JSON APIs, not server-side rendered applications.
*   **Secure by Default:** Sensible defaults for security headers, JWT, CSRF, and rate limiting are provided out-of-the-box.
*   **FastAPI-inspired DX:** Aims for an intuitive and expressive API, inspired by the clarity and developer experience of frameworks like FastAPI.
*   **Go-native Routing & Middleware:** Built on Go's standard library `net/http`, ensuring minimal overhead and excellent performance.
*   **Explicit, No Magic:** Clear and explicit middleware application gives you full control over the request processing pipeline.

## 🚀 Quickstart

### Installation

To add `goBastion Lib` to your Go project, use the `go get` command:

```bash
//external
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
bastion
go get github.com/alejandrombjs/go-bastion-lib
go mod tidy
```

### Minimal API Example

Create a `main.go` file:

```go
package main

import (
	"log"
	"time"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/response"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
	"github.com/alejandrombjs/go-bastion-lib/pkg/security"
)

func main() {
	cfg := bastion.DefaultConfig()
	cfg.EnableJWT = true
	cfg.JWTSecret = "super-secret-jwt-key-that-is-at-least-32-bytes-long" // IMPORTANT: Change in production!
	cfg.EnableRateLimit = true
	cfg.EnableSecurityHeaders = true
	cfg.Port = 8080

	app := bastion.NewApp(cfg)

	app.Use(
		middleware.RequestID(),
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
		middleware.SecurityHeaders(middleware.DefaultSecurityHeaders()),
	)

	if cfg.EnableRateLimit {
		app.Use(middleware.RateLimit(cfg.RateLimitRequests, cfg.RateLimitWindow))
	}

	r := app.Router()

	r.GET("/api/health", func(ctx *router.Context) {
		response.JSON(ctx, 200, map[string]string{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	authGroup := r.Group("/api/auth")
	authGroup.POST("/login", func(ctx *router.Context) {
		var loginReq struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := ctx.BindJSON(&loginReq); err != nil {
			response.Error(ctx, 400, "invalid_request", "Invalid JSON body")
			return
		}
		if loginReq.Username == "testuser" && loginReq.Password == "testpass" {
			token, err := security.GenerateAccessToken(
				loginReq.Username,
				15*time.Minute,
				cfg.JWTSecret,
				map[string]any{"role": "user"},
			)
			if err != nil {
				response.Error(ctx, 500, "internal_error", "Failed to generate token")
				return
			}
			response.JSON(ctx, 200, map[string]any{
				"access_token": token,
				"token_type":   "bearer",
				"expires_in":   900,
			})
		} else {
			response.Error(ctx, 401, "unauthorized", "Invalid credentials")
		}
	})

	protectedApiGroup := r.Group("/api")
	protectedApiGroup.Use(middleware.JWTAuth(cfg.JWTSecret))

	protectedApiGroup.GET("/users", func(ctx *router.Context) {
		claims, ok := ctx.Get("userClaims")
		if !ok {
			response.Error(ctx, 401, "unauthorized", "User claims not found")
			return
		}
		response.Success(ctx, 200, map[string]any{
			"message": "Welcome to protected user data!",
			"user":    claims,
			"data":    []string{"item1", "item2"},
		})
	})

	log.Printf("Server starting on port %d", cfg.Port)
	if err := app.RunWithGracefulShutdown(); err != nil {
		log.Fatalf("Server failed to start or shut down gracefully: %v", err)
	}
}
```

Run the application:

```bash
go run main.go
```

Then test with `curl`:

```bash
curl http://localhost:8080/api/health
# Expected: {"status":"ok","timestamp":"..."}
```

## 📚 Documentation

For a comprehensive guide on `goBastion Lib`'s features, architecture, security, and API reference, please visit our full documentation site:

👉 **[Open Documentation Site](docs/index.html)**

## ✨ Examples

Explore a curated collection of runnable examples demonstrating various features and integration patterns with `goBastion Lib`. Each example is a self-contained mini-application designed to showcase a specific use case, built with `goBastion Lib` and focusing exclusively on JSON-based backend APIs.

| Example                                  | Description                                                              |
| :--------------------------------------- | :----------------------------------------------------------------------- |
| **Hello World + 404 Handler**            | Minimal API + custom JSON 404 responses.                                 |
| **JWT Auth (Core)**                      | Simple login + protected `/profile` route.                               |
| **Auth + JWT + Docker + Postgres**       | Auth backed by Postgres and Docker Compose.                              |
| **CSRF Protection**                      | Demonstrates CSRF middleware integration.                                |
| **GORM + PostgreSQL CRUD**               | Full users CRUD with GORM & PostgreSQL.                                  |
| **File Upload**                          | Multipart file upload and safe storage.                                  |
| **WebSocket**                            | Real-time echo/chat over WebSockets.                                     |
| **OpenAPI + Swagger**                    | Generating and serving OpenAPI/Swagger docs.                             |
| **Testing (Unit & Integration)**         | Patterns for testing `goBastion Lib` handlers.                           |
| **Docker + Nginx Reverse Proxy**         | Reverse proxy setup for production-like deploy.                          |
| **Graceful Shutdown**                    | Clean SIGINT/SIGTERM handling and shutdown.                              |

👉 For a richer overview, detailed descriptions, and interactive filtering of examples, please see the **[Examples Catalog](docs/index.html#examples-catalog)** section on our documentation site.

## 🧪 Automated API Tests

This project includes an automated API testing script (`scripts/test-api.sh`) that uses `curl` and `jq` to perform integration tests against a running instance of the `goBastion Lib` application.

To run the tests:

```bash
./scripts/test-api.sh
```

This script will:

1.  Build and start the `examples/basic/main.go` server in the background.
2.  Wait for the server to become healthy.
3.  Execute a series of `curl` commands to test various endpoints (e.g., health checks, protected routes, error handling).
4.  Validate JSON responses using `jq`.
5.  Report test results.
6.  Gracefully shut down the server.

## 🤝 Contributing

We welcome contributions! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

`goBastion Lib` is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
