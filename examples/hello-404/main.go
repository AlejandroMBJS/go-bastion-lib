package main

import (
	"log"
	"net/http"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

func main() {
	cfg := bastion.DefaultConfig()
	cfg.Port = 8080 // Set a specific port for this example

	app := bastion.NewApp(cfg)
	r := app.Router()

	// --- Custom 404 Not Found Handler ---
	// This demonstrates how to override the default 404 JSON response.
	r.SetNotFoundHandler(func(ctx *router.Context) {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"error":   "not_found",
			"message": "The requested resource could not be found on this server.",
		})
	})

	// --- Hello World Endpoint ---
	r.GET("/", func(ctx *router.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "Hello, World!",
		})
	})

	log.Printf("Hello-404 example server starting on :%d", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
