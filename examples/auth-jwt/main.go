package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/response"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
	"github.com/alejandrombjs/go-bastion-lib/pkg/security"
)

// User represents a simple user model for this example
type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // In real app, store hashed password
	Role     string `json:"role"`
}

// Hardcoded user for demonstration purposes
var exampleUser = User{
	Username: "testuser",
	Password: "testpass", // In real app, this would be a hashed password
	Role:     "user",
}

func main() {
	cfg := bastion.DefaultConfig()
	cfg.Port = 8081 // Use a different port for this example
	cfg.EnableJWT = true
	// IMPORTANT: In a real application, load this from environment variables!
	cfg.JWTSecret = "super-secret-jwt-key-that-is-at-least-32-bytes-long-for-auth-jwt-example"

	app := bastion.NewApp(cfg)
	r := app.Router()

	// Global middlewares
	app.Use(
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
	)

	// --- Public Routes ---
	// Login endpoint
	r.POST("/auth/login", func(ctx *router.Context) {
		var loginReq struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := ctx.BindJSON(&loginReq); err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
			return
		}

		// Authenticate user (in real app, check hashed password from DB)
		if loginReq.Username == exampleUser.Username && loginReq.Password == exampleUser.Password {
			// Generate JWT token
			token, err := security.GenerateAccessToken(
				exampleUser.Username,
				15*time.Minute, // Access token valid for 15 minutes
				cfg.JWTSecret,
				map[string]any{"role": exampleUser.Role}, // Custom claims
			)
			if err != nil {
				response.Error(ctx, http.StatusInternalServerError, "internal_error", "Failed to generate token")
				return
			}

			response.JSON(ctx, http.StatusOK, map[string]any{
				"access_token": token,
				"token_type":   "bearer",
				"expires_in":   900, // 15 minutes in seconds
			})
		} else {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized", "Invalid credentials")
		}
	})

	// --- Protected Routes ---
	// Create a group for protected routes and apply JWT middleware
	protected := r.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))

	protected.GET("/profile", func(ctx *router.Context) {
		// Access user claims from context (set by JWTAuth middleware)
		claims, ok := ctx.Get("userClaims").(*security.Claims)
		if !ok || claims == nil {
			// This should ideally not happen if JWTAuth middleware works correctly
			response.Error(ctx, http.StatusUnauthorized, "unauthorized", "User claims not found")
			return
		}

		response.Success(ctx, http.StatusOK, map[string]any{
			"message":   "Welcome to your profile!",
			"username":  claims.Subject,
			"user_role": claims.Extra["role"],
			"claims":    claims,
		})
	})

	log.Printf("Auth-JWT example server starting on :%d", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
