package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

func main() {
	cfg := bastion.DefaultConfig()
	cfg.Port = 8089 // Use a different port for this example

	app := bastion.NewApp(cfg)
	r := app.Router()

	// Global middlewares
	app.Use(
		middleware.RequestID(),
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
	)

	// --- Long-running Endpoint ---
	// This endpoint simulates a task that takes some time to complete.
	// It will be used to demonstrate graceful shutdown.
	r.GET("/long-task", func(ctx *router.Context) {
		log.Printf("Request %s: Starting long task...", ctx.GetString("requestID"))
		time.Sleep(5 * time.Second) // Simulate work
		log.Printf("Request %s: Long task finished.", ctx.GetString("requestID"))
		ctx.JSON(http.StatusOK, map[string]string{"message": "Long task completed"})
	})

	// --- Health Check ---
	r.GET("/health", func(ctx *router.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Start server in a goroutine
	go func() {
		log.Printf("Graceful Shutdown example server starting on :%d", cfg.Port)
		if err := app.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with a timeout for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
}
