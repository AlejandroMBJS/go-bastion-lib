// Package response provides helper functions for HTTP responses.
package response

import (
	"net/http" // Added for http.StatusText and http.StatusInternalServerError

	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
	"github.com/alejandrombjs/go-bastion-lib/pkg/templating" // Added templating import
)

// JSON sends a JSON response with the given status code.
func JSON(ctx *router.Context, status int, payload any) {
	ctx.JSON(status, payload)
}

// Error sends an error response with the given status code and error details.
func Error(ctx *router.Context, status int, code, message string) {
	ctx.JSON(status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

// Success sends a success response with the given status code and data.
func Success(ctx *router.Context, status int, data any) {
	ctx.JSON(status, map[string]any{
		"data": data,
	})
}

// Created sends a 201 Created response with location header.
func Created(ctx *router.Context, location string, data any) {
	ctx.ResponseWriter().Header().Set("Location", location)
	Success(ctx, 201, data)
}

// NoContent sends a 204 No Content response.
func NoContent(ctx *router.Context) {
	ctx.Status(204)
}

// H is a convenient alias for templating.H, allowing user code to use response.H{...}.
type H = templating.H

// HTML renders an HTML template with the given status code and data.
func HTML(ctx *router.Context, status int, name string, data any) {
	w := ctx.ResponseWriter() // Corrected: Use ResponseWriter() method

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if err := templating.Render(w, name, data); err != nil {
		// TODO: integrate with your logging if you have a logger available.
		// For now, send a generic 500.
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
