package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

// CSRFConfig holds CSRF middleware configuration.
type CSRFConfig struct {
	CookieName string
	HeaderName string
	Secure     bool
	SameSite   http.SameSite
}

// DefaultCSRFConfig returns default CSRF configuration.
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		CookieName: "csrf_token",
		HeaderName: "X-CSRF-Token",
		Secure:     true,
		SameSite:   http.SameSiteStrictMode,
	}
}

// CSRFMiddleware creates a CSRF protection middleware.
func CSRFMiddleware(cfg CSRFConfig) router.Middleware {
	return func(next router.Handler) router.Handler {
		return func(ctx *router.Context) {
			// Skip CSRF for safe methods
			if isSafeMethod(ctx.Request().Method) {
				// Still set CSRF cookie for safe methods
				setCSRFCookie(ctx, cfg)
				next(ctx)
				return
			}

			// For unsafe methods, verify CSRF token
			cookieToken, err := ctx.Request().Cookie(cfg.CookieName)
			if err != nil {
				ctx.JSON(http.StatusForbidden, map[string]string{
					"error": "csrf_failed",
				})
				return
			}

			// Get token from header or form
			headerToken := ctx.Request().Header.Get(cfg.HeaderName)
			if headerToken == "" {
				headerToken = ctx.FormValue(cfg.HeaderName)
			}
			// Verify tokens match
			if cookieToken.Value == "" || cookieToken.Value != headerToken {
				ctx.JSON(http.StatusForbidden, map[string]string{
					"error": "csrf_failed",
				})
				return
			}

			// Generate new token for next request
			setCSRFCookie(ctx, cfg)
			next(ctx)
		}
	}
}

// isSafeMethod checks if the HTTP method is considered safe (doesn't change state).
func isSafeMethod(method string) bool {
	safeMethods := map[string]bool{
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodOptions: true,
	}
	return safeMethods[method]
}

// setCSRFCookie sets or renews the CSRF cookie.
func setCSRFCookie(ctx *router.Context, cfg CSRFConfig) {
	token, err := generateCSRFToken()
	if err != nil {
		return
	}

	cookie := &http.Cookie{
		Name:     cfg.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: false, // Must be accessible by JavaScript
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
	}

	http.SetCookie(ctx.ResponseWriter(), cookie)
}

// generateCSRFToken generates a random CSRF token.
func generateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
