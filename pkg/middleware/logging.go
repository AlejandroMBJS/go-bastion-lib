package middleware

import (
	"log"
	"time"

	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

// Logging creates a logging middleware.
func Logging(logger *log.Logger) router.Middleware {
	return func(next router.Handler) router.Handler {
		return func(ctx *router.Context) {
			start := time.Now()

			// Get request ID
			requestID, _ := ctx.GetString("requestID")

			// Call next handler
			next(ctx)

			// Calculate duration
			duration := time.Since(start)

			// Log the request
			logger.Printf("[%s] %s %s %d %v",
				requestID,
				ctx.Request().Method,
				ctx.Request().URL.Path,
				ctx.StatusCode(),
				duration,
			)
		}
	}
}

// DefaultLogging uses the standard logger.
func DefaultLogging() router.Middleware {
	return Logging(log.Default())
}
