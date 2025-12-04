// tests/router_test.go
package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

func TestRouterBasic(t *testing.T) {
	r := router.New()

	r.GET("/test", func(ctx *router.Context) {
		ctx.JSON(200, map[string]string{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.Handler().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRouterParams(t *testing.T) {
	r := router.New()

	r.GET("/users/:id", func(ctx *router.Context) {
		id := ctx.Param("id")
		ctx.JSON(200, map[string]string{"id": id})
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	r.Handler().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRouterGroup(t *testing.T) {
	r := router.New()

	api := r.Group("/api")
	api.GET("/users", func(ctx *router.Context) {
		ctx.JSON(200, map[string]string{"message": "users"})
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()

	r.Handler().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
