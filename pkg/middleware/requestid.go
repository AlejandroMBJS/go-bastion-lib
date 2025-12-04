package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

// RequestID generates a unique request ID for each request.
func RequestID() router.Middleware {
	return func(next router.Handler) router.Handler {
		return func(ctx *router.Context) {
			// Generate random request ID
			id := generateRequestID()

			// Set in context
			ctx.Set("requestID", id)

			// Set in response header
			ctx.ResponseWriter().Header().Set("X-Request-ID", id)

			next(ctx)
		}
	}
}

// generateRequestID generates a random 32-character hex string.
func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
