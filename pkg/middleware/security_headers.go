package middleware

import (
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

// SecurityHeaders adds security headers to responses.
func SecurityHeaders(cfg SecurityHeadersConfig) router.Middleware {
	return func(next router.Handler) router.Handler {
		return func(ctx *router.Context) {
			w := ctx.ResponseWriter()

			// Set security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "no-referrer")

			if cfg.CSP != "" {
				w.Header().Set("Content-Security-Policy", cfg.CSP)
			}

			if cfg.HSTS != "" {
				w.Header().Set("Strict-Transport-Security", cfg.HSTS)
			}

			next(ctx)
		}
	}
}

// SecurityHeadersConfig holds security header configuration.
type SecurityHeadersConfig struct {
	CSP  string // Content-Security-Policy
	HSTS string // Strict-Transport-Security
}

// DefaultSecurityHeaders returns a SecurityHeadersConfig with safe defaults.
func DefaultSecurityHeaders() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		CSP:  "default-src 'self';",
		HSTS: "max-age=31536000; includeSubDomains",
	}
}
