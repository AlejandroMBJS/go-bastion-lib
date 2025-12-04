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
	// Load configuration
	cfg := bastion.DefaultConfig()
	cfg.EnableJWT = true
	cfg.JWTSecret = "your-256-bit-secret-key-here-change-in-production"
	cfg.EnableRateLimit = true
	cfg.EnableSecurityHeaders = true

	// Create app
	app := bastion.NewApp(cfg)

	// Register global middlewares
	app.Use(middleware.RequestID())
	app.Use(middleware.DefaultLogging())
	app.Use(middleware.DefaultRecovery())
	app.Use(middleware.SecurityHeaders(middleware.DefaultSecurityHeaders()))

	// Apply rate limiting if enabled
	if cfg.EnableRateLimit {
		app.Use(middleware.RateLimit(cfg.RateLimitRequests, cfg.RateLimitWindow))
	}

	// Get router
	r := app.Router()

	// Public routes
	r.GET("/api/health", healthHandler)

	// Auth routes
	auth := r.Group("/api/auth")
	auth.POST("/login", loginHandler)

	// Protected routes (require JWT)
	api := r.Group("/api")
	api.Use(middleware.JWTAuth(cfg.JWTSecret))

	users := api.Group("/users")
	users.GET("/", listUsersHandler)
	users.POST("/", createUserHandler)

	// Start server
	if err := app.RunWithGracefulShutdown(); err != nil {
		log.Fatal(err)
	}
}

// healthHandler returns service health status
func healthHandler(ctx *router.Context) {
	response.JSON(ctx, 200, map[string]string{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// loginHandler handles user authentication
func loginHandler(ctx *router.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		response.Error(ctx, 400, "invalid_request", "Invalid JSON body")
		return
	}

	// In a real application, you would:
	// 1. Look up user in database
	// 2. Verify password hash
	// 3. Generate appropriate claims

	// For this example, accept any credentials
	// WARNING: Never do this in production!

	// Generate JWT token
	token, err := security.GenerateAccessToken(
		req.Username,
		15*time.Minute,
		"your-256-bit-secret-key-here-change-in-production",
		map[string]any{
			"username": req.Username,
			"role":     "user",
		},
	)

	if err != nil {
		response.Error(ctx, 500, "internal_error", "Failed to generate token")
		return
	}

	response.JSON(ctx, 200, map[string]any{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   900, // 15 minutes in seconds
	})
}

// listUsersHandler returns a list of users (protected)
func listUsersHandler(ctx *router.Context) {
	// Get user claims from context
	claims, ok := ctx.Get("userClaims")
	if !ok {
		response.Error(ctx, 401, "unauthorized", "User claims not found")
		return
	}

	// In a real application, you would fetch from database
	// For this example, return dummy data
	users := []map[string]any{
		{
			"id":       1,
			"username": "john_doe",
			"email":    "john@example.com",
			"role":     "user",
		},
		{
			"id":       2,
			"username": "jane_doe",
			"email":    "jane@example.com",
			"role":     "admin",
		},
	}

	response.Success(ctx, 200, map[string]any{
		"users":     users,
		"count":     len(users),
		"requester": claims,
	})
}

// createUserHandler creates a new user (protected)
func createUserHandler(ctx *router.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		response.Error(ctx, 400, "invalid_request", "Invalid JSON body")
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		response.Error(ctx, 400, "validation_error", "All fields are required")
		return
	}

	// Hash password (in real app, do more validation)
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		response.Error(ctx, 500, "internal_error", "Failed to hash password")
		return
	}

	// In a real application, you would:
	// 1. Check if user already exists
	// 2. Create user in database
	// 3. Send welcome email, etc.

	// For this example, return success
	response.Created(ctx, "/api/users/3", map[string]any{
		"id":       3,
		"username": req.Username,
		"email":    req.Email,
		"message":  "User created successfully",
		"note":     "Password hashed (not stored in response): " + hashedPassword[:10] + "...",
	})
}
