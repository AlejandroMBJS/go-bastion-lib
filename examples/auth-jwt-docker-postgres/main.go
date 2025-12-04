package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/response"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
	"github.com/alejandrombjs/go-bastion-lib/pkg/security"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User model for GORM
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // Store hashed password, hide from JSON
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Role     string `gorm:"default:'user'" json:"role"`
}

// LoginRequest struct for binding JSON
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest struct for binding JSON
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func main() {
	// Load configuration from environment variables
	cfg := bastion.LoadConfigFromEnv()
	cfg.Port = 8082 // Use a different port for this example
	cfg.EnableJWT = true
	// JWT_SECRET will be loaded from environment by LoadConfigFromEnv
	// If not set, it will be empty, causing token generation to fail.

	// --- Database Setup ---
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var db *gorm.DB
	var err error
	// Retry connecting to DB as it might not be ready immediately in Docker Compose
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err == nil {
			log.Println("Successfully connected to the database")
			break
		}
		log.Printf("Failed to connect to database, retrying in 2 seconds... (%d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after multiple retries: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("Database migration complete.")

	// Seed a default admin user if not exists
	var adminUser User
	if db.Where("username = ?", "admin").First(&adminUser).Error == gorm.ErrRecordNotFound {
		hashedPassword, _ := security.HashPassword("adminpass")
		adminUser = User{
			Username: "admin",
			Password: hashedPassword,
			Email:    "admin@example.com",
			Role:     "admin",
		}
		db.Create(&adminUser)
		log.Println("Default admin user created.")
	}

	// --- Bastion App Setup ---
	app := bastion.NewApp(cfg)
	r := app.Router()

	// Global middlewares
	app.Use(
		middleware.RequestID(),
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
	)

	// --- Public Routes ---
	// User registration
	r.POST("/auth/register", func(ctx *router.Context) {
		var req RegisterRequest
		if err := ctx.BindJSON(&req); err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
			return
		}

		if req.Username == "" || req.Password == "" || req.Email == "" {
			response.Error(ctx, http.StatusBadRequest, "validation_error", "Username, password, and email are required")
			return
		}

		// Check if username or email already exists
		var existingUser User
		if db.Where("username = ?", req.Username).Or("email = ?", req.Email).First(&existingUser).Error == nil {
			response.Error(ctx, http.StatusConflict, "user_exists", "Username or email already registered")
			return
		}

		hashedPassword, err := security.HashPassword(req.Password)
		if err != nil {
			response.Error(ctx, http.StatusInternalServerError, "internal_error", "Failed to hash password")
			return
		}

		user := User{
			Username: req.Username,
			Password: hashedPassword,
			Email:    req.Email,
			Role:     "user", // Default role
		}

		if result := db.Create(&user); result.Error != nil {
			response.Error(ctx, http.StatusInternalServerError, "db_error", "Failed to create user")
			return
		}

		response.Created(ctx, fmt.Sprintf("/users/%d", user.ID), map[string]string{
			"message":  "User registered successfully",
			"username": user.Username,
		})
	})

	// Login endpoint
	r.POST("/auth/login", func(ctx *router.Context) {
		var req LoginRequest
		if err := ctx.BindJSON(&req); err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
			return
		}

		var user User
		if result := db.Where("username = ?", req.Username).First(&user); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				response.Error(ctx, http.StatusUnauthorized, "unauthorized", "Invalid credentials")
				return
			}
			response.Error(ctx, http.StatusInternalServerError, "db_error", "Database query failed")
			return
		}

		if !security.CheckPasswordHash(req.Password, user.Password) {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized", "Invalid credentials")
			return
		}

		// Generate JWT token
		token, err := security.GenerateAccessToken(
			user.Username,
			15*time.Minute, // Access token valid for 15 minutes
			cfg.JWTSecret,
			map[string]any{"role": user.Role, "email": user.Email}, // Custom claims
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
	})

	// --- Protected Routes ---
	protected := r.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))

	protected.GET("/profile", func(ctx *router.Context) {
		claims, ok := ctx.Get("userClaims").(*security.Claims)
		if !ok || claims == nil {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized", "User claims not found")
			return
		}

		response.Success(ctx, http.StatusOK, map[string]any{
			"message":   "Welcome to your profile!",
			"username":  claims.Subject,
			"user_role": claims.Extra["role"],
			"user_email": claims.Extra["email"],
			"claims":    claims,
		})
	})

	// Admin-only route example
	protected.GET("/admin/dashboard", func(ctx *router.Context) {
		claims, ok := ctx.Get("userClaims").(*security.Claims)
		if !ok || claims == nil {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized", "User claims not found")
			return
		}

		if role, ok := claims.Extra["role"].(string); !ok || role != "admin" {
			response.Error(ctx, http.StatusForbidden, "forbidden", "Admin access required")
			return
		}

		response.Success(ctx, http.StatusOK, map[string]string{
			"message": "Welcome to the admin dashboard!",
			"admin":   claims.Subject,
		})
	})

	log.Printf("Auth-JWT-Docker-Postgres example server starting on :%d", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
