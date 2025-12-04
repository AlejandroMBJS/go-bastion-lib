// Package bastion provides the main application framework.
package bastion

import (
	"context"
	"fmt"
	"log" // Added for logging fatal errors
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
	"github.com/alejandrombjs/go-bastion-lib/pkg/templating" // Added templating import
)

// App represents the main application instance.
type App struct {
	config   Config
	server   *http.Server
	router   *router.Router
	shutdown chan os.Signal
}

// NewApp creates a new application instance with the given configuration.
func NewApp(cfg Config) *App {
	r := router.New()

	app := &App{
		config:   cfg,
		router:   r,
		shutdown: make(chan os.Signal, 1),
	}

	// Initialize the default templating engine
	isDevelopment := cfg.Env == "development"
	err := templating.InitDefault(templating.Options{
		Root:         cfg.TemplateRoot, // Use TemplateRoot from config
		Extensions:   []string{".gb.html", ".html"},
		CacheEnabled: !isDevelopment,
		Debug:        isDevelopment,
		Funcs:        nil,
	})
	if err != nil {
		log.Fatalf("failed to initialize templating engine: %v", err)
	}

	// Create HTTP server with configured timeouts
	app.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r.Handler(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return app
}

// Router returns the root router of the application.
// This allows users to register routes and middlewares.
func (a *App) Router() *router.Router {
	return a.router
}

// Use registers global middleware that will be applied to all routes.
func (a *App) Use(middleware ...router.Middleware) {
	a.router.Use(middleware...)
}

// Run starts the HTTP server and listens for incoming requests.
// It blocks until the server is stopped.
func (a *App) Run() error {
	fmt.Printf("Starting server on :%d\n", a.config.Port)
	fmt.Printf("Environment: %s\n", a.config.Env)

	// Start server in a goroutine so we can handle shutdown
	errChan := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Listen for shutdown signals
	signal.Notify(a.shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-a.shutdown:
		fmt.Printf("Received signal: %v\n", sig)
	case err := <-errChan:
		return err
	}

	return nil
}

// Shutdown gracefully shuts down the server with the given timeout.
func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

// RunWithGracefulShutdown starts the server and handles graceful shutdown.
func (a *App) RunWithGracefulShutdown() error {
	errChan := make(chan error, 1)

	// Start server
	go func() {
		fmt.Printf("Starting server on :%d\n", a.config.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for interrupt signal
	signal.Notify(a.shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-a.shutdown:
		fmt.Printf("Received signal: %v, shutting down gracefully...\n", sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := a.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %v", err)
		}

		fmt.Println("Server stopped gracefully")
		return nil

	case err := <-errChan:
		return err
	}
}
