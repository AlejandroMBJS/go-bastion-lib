package main

import (
	"log"
	"net/http"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/response"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

func main() {
	cfg := bastion.DefaultConfig()
	cfg.Port = 8083 // Use a different port for this example
	cfg.EnableCSRF = true // Enable CSRF protection

	app := bastion.NewApp(cfg)
	r := app.Router()

	// Global middlewares
	app.Use(
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
	)

	// Apply CSRF middleware globally
	// For this example, we'll use the default config.
	// In a real app, ensure Secure is true for HTTPS.
	app.Use(middleware.CSRFMiddleware(middleware.DefaultCSRFConfig()))

	// --- Public Route to get CSRF token ---
	// This simulates a frontend requesting a page that sets the CSRF cookie.
	// For simplicity, we'll just return a message. The CSRF cookie will be set
	// by the middleware on any GET request.
	r.GET("/form", func(ctx *router.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "Load this form to get a CSRF cookie. Then extract X-CSRF-Token from the cookie.",
			"note":    "Check your browser's cookies for 'csrf_token' after this request.",
		})
	})

	// --- Protected Route requiring CSRF token ---
	r.POST("/form/submit", func(ctx *router.Context) {
		var data struct {
			Message string `json:"message"`
		}
		if err := ctx.BindJSON(&data); err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
			return
		}

		ctx.JSON(http.StatusOK, map[string]string{
			"status":  "success",
			"message": "Form submitted successfully!",
			"data":    data.Message,
		})
	})

	log.Printf("CSRF example server starting on :%d", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
