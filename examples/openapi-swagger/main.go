package main

import (
	"log"
	"net/http"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

func main() {
	cfg := bastion.DefaultConfig()
	cfg.Port = 8087 // Use a different port for this example

	app := bastion.NewApp(cfg)
	r := app.Router()

	// Global middlewares
	app.Use(
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
	)

	// --- Serve OpenAPI Specification ---
	// This endpoint serves the openapi.yaml file directly.
	r.GET("/openapi.yaml", func(ctx *router.Context) {
		http.ServeFile(ctx.ResponseWriter(), ctx.Request().Request, "./openapi.yaml")
	})

	// --- Serve Swagger UI ---
	// This endpoint redirects to a public Swagger UI instance,
	// configured to load our local openapi.yaml.
	// In a production setup, you would typically serve Swagger UI static files
	// directly from your application or a CDN.
	r.GET("/docs", func(ctx *router.Context) {
		swaggerUIURL := fmt.Sprintf("https://petstore.swagger.io/?url=http://localhost:%d/openapi.yaml", cfg.Port)
		http.Redirect(ctx.ResponseWriter(), ctx.Request().Request, swaggerUIURL, http.StatusMovedPermanently)
	})

	// --- Example API Endpoint (documented in openapi.yaml) ---
	r.GET("/api/hello", func(ctx *router.Context) {
		name := ctx.Query("name")
		if name == "" {
			name = "World"
		}
		ctx.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Hello, %s!", name)})
	})

	log.Printf("OpenAPI-Swagger example server starting on :%d", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
