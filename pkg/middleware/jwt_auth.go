package middleware

import (
	"net/http"
	"strings"

	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
	"github.com/alejandrombjs/go-bastion-lib/pkg/security"
)

// JWTAuth creates a JWT authentication middleware.
func JWTAuth(secret string) router.Middleware {
	return func(next router.Handler) router.Handler {
		return func(ctx *router.Context) {
			// Extract token from Authorization header
			authHeader := ctx.Request().Header.Get("Authorization")
			if authHeader == "" {
				ctx.JSON(http.StatusUnauthorized, map[string]string{
					"error": "unauthorized",
				})
				return
			}
			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				ctx.JSON(http.StatusUnauthorized, map[string]string{
					"error": "unauthorized",
				})
				return
			}

			token := parts[1]

			// Parse and validate token
			claims, err := security.ParseAndValidateToken(token, secret)
			if err != nil {
				if err == security.ErrExpiredToken {
					ctx.JSON(http.StatusUnauthorized, map[string]string{
						"error": "token_expired",
					})
				} else {
					ctx.JSON(http.StatusUnauthorized, map[string]string{
						"error": "unauthorized",
					})
				}
				return
			}

			// Store claims in context
			ctx.Set("userClaims", claims)

			next(ctx)
		}
	}
}
