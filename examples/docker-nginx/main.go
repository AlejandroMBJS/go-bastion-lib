package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

func main() {
	cfg := bastion.DefaultConfig()
	// Port will be set by environment variable in Docker Compose
	// Default to 8080 if not set, but Dockerfile will expose 8080
	if os.Getenv("PORT") != "" {
		p, err := strconv.Atoi(os.Getenv("PORT"))
		if err == nil {
			cfg.Port = p
		}
	} else {
		cfg.Port = 8080
	}

	app := bastion.NewApp(cfg)
	r := app.Router()

	// Global middlewares
	app.Use(
		middleware.RequestID(),
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
	)

	// --- API Endpoints ---
	r.GET("/api/hello", func(ctx *router.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "Hello from Go Bastion App!",
			"served_by": fmt.Sprintf("Port %d", cfg.Port),
		})
	})

	r.GET("/api/headers", func(ctx *router.Context) {
		headers := make(map[string]string)
		for name, values := range ctx.Request().Header {
			headers[name] = strings.Join(values, ", ")
		}
		ctx.JSON(http.StatusOK, map[string]any{
			"message": "Request headers received by Go Bastion App",
			"headers": headers,
		})
	})

	log.Printf("Go Bastion App starting on :%d", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
