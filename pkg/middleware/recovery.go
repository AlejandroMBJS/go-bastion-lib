package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

// Recovery creates a panic recovery middleware.
func Recovery(logger *log.Logger) router.Middleware {
	return func(next router.Handler) router.Handler {
		return func(ctx *router.Context) {
			defer func() {
				if r := recover(); r != nil {
					// Log the panic
					logger.Printf("panic: %v\n%s", r, debug.Stack())

					// Get request ID
					requestID, _ := ctx.GetString("requestID")

					// Send error response
					ctx.JSON(http.StatusInternalServerError, map[string]any{
						"error":      "internal_server_error",
						"request_id": requestID,
					})
				}
			}()

			next(ctx)
		}
	}
}

// DefaultRecovery uses the standard logger.
func DefaultRecovery() router.Middleware {
	return Recovery(log.Default())
}
