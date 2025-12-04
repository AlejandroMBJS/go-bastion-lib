package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/response"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Product model for GORM
type Product struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	Stock       int     `gorm:"not null" json:"stock"`
}

// CreateProductRequest struct for binding JSON
type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// UpdateProductRequest struct for binding JSON
type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
}

func main() {
	// Load configuration from environment variables
	cfg := bastion.LoadConfigFromEnv()
	cfg.Port = 8084 // Use a different port for this example

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
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("Database migration complete.")

	// --- Bastion App Setup ---
	app := bastion.NewApp(cfg)
	r := app.Router()

	// Global middlewares
	app.Use(
		middleware.RequestID(),
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
	)

	// --- Product CRUD Routes ---
	productsGroup := r.Group("/products")

	// GET /products - List all products
	productsGroup.GET("/", func(ctx *router.Context) {
		var products []Product
		if result := db.Find(&products); result.Error != nil {
			response.Error(ctx, http.StatusInternalServerError, "db_error", "Failed to retrieve products")
			return
		}
		response.Success(ctx, http.StatusOK, products)
	})

	// POST /products - Create a new product
	productsGroup.POST("/", func(ctx *router.Context) {
		var req CreateProductRequest
		if err := ctx.BindJSON(&req); err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
			return
		}

		if req.Name == "" || req.Price <= 0 {
			response.Error(ctx, http.StatusBadRequest, "validation_error", "Product name and price must be provided and valid")
			return
		}

		product := Product{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
		}

		if result := db.Create(&product); result.Error != nil {
			response.Error(ctx, http.StatusInternalServerError, "db_error", "Failed to create product")
			return
		}

		response.Created(ctx, fmt.Sprintf("/products/%d", product.ID), product)
	})

	// GET /products/:id - Get a single product by ID
	productsGroup.GET("/:id", func(ctx *router.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_id", "Invalid product ID")
			return
		}

		var product Product
		if result := db.First(&product, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				response.Error(ctx, http.StatusNotFound, "not_found", "Product not found")
				return
			}
			response.Error(ctx, http.StatusInternalServerError, "db_error", "Failed to retrieve product")
			return
		}
		response.Success(ctx, http.StatusOK, product)
	})

	// PUT /products/:id - Update a product by ID
	productsGroup.PUT("/:id", func(ctx *router.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_id", "Invalid product ID")
			return
		}

		var req UpdateProductRequest
		if err := ctx.BindJSON(&req); err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
			return
		}

		var product Product
		if result := db.First(&product, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				response.Error(ctx, http.StatusNotFound, "not_found", "Product not found")
				return
			}
			response.Error(ctx, http.StatusInternalServerError, "db_error", "Failed to retrieve product for update")
			return
		}

		// Apply updates only if fields are provided
		if req.Name != nil {
			product.Name = *req.Name
		}
		if req.Description != nil {
			product.Description = *req.Description
		}
		if req.Price != nil {
			product.Price = *req.Price
		}
		if req.Stock != nil {
			product.Stock = *req.Stock
		}

		if result := db.Save(&product); result.Error != nil {
			response.Error(ctx, http.StatusInternalServerError, "db_error", "Failed to update product")
			return
		}

		response.Success(ctx, http.StatusOK, product)
	})

	// DELETE /products/:id - Delete a product by ID
	productsGroup.DELETE("/:id", func(ctx *router.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_id", "Invalid product ID")
			return
		}

		if result := db.Delete(&Product{}, id); result.Error != nil {
			response.Error(ctx, http.StatusInternalServerError, "db_error", "Failed to delete product")
			return
		}

		response.NoContent(ctx) // 204 No Content for successful deletion
	})

	log.Printf("GORM-PostgreSQL-CRUD example server starting on :%d", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
