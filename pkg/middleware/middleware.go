// Package middleware provides common HTTP middlewares.
package middleware

import (
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

// Chain creates a middleware chain.
func Chain(middlewares ...router.Middleware) router.Middleware {
	return func(next router.Handler) router.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
