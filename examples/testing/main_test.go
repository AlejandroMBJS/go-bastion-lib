package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
	"github.com/stretchr/testify/assert"
)

// --- Unit Tests for Handlers ---

func TestPingHandler(t *testing.T) {
	// Create a mock Context
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ping", nil)
	ctx := router.NewContext(w, r)

	PingHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "pong", resp["message"])
}

func TestListUsersHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	ctx := router.NewContext(w, r)

	ListUsersHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data []User `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Data, 2) // Assuming initial users are Alice and Bob
}

func TestCreateUserHandler_Success(t *testing.T) {
	initialNextUserID := nextUserID // Save initial state
	initialUsers := make(map[int]User)
	for k, v := range users {
		initialUsers[k] = v
	}

	defer func() { // Restore initial state after test
		nextUserID = initialNextUserID
		users = initialUsers
	}()

	body := []byte(`{"name":"Charlie","age":25}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := router.NewContext(w, r)

	CreateUserHandler(ctx)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp User
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Charlie", resp.Name)
	assert.Equal(t, 25, resp.Age)
	assert.Equal(t, initialNextUserID, resp.ID) // Check if ID is correctly assigned
	assert.Contains(t, users, initialNextUserID) // Check if user is added to map
}

func TestCreateUserHandler_InvalidInput(t *testing.T) {
	body := []byte(`{"name":"","age":-5}`) // Invalid name and age
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := router.NewContext(w, r)

	CreateUserHandler(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", resp["error"]["code"])
}

func TestGetUserHandler_Success(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	ctx := router.NewContext(w, r)
	ctx.SetParam("id", "1") // Manually set path parameter

	GetUserHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data User `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", resp.Data.Name)
}

func TestGetUserHandler_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	ctx := router.NewContext(w, r)
	ctx.SetParam("id", "99") // Manually set path parameter

	GetUserHandler(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", resp["error"]["code"])
}

// --- Integration Tests for Router ---

func TestRouterIntegration(t *testing.T) {
	// Setup the full application and router for integration testing
	cfg := bastion.DefaultConfig()
	app := bastion.NewApp(cfg)
	r := app.Router()

	// Apply middlewares (important for integration tests)
	r.Use(middleware.DefaultRecovery()) // Ensure recovery works

	// Register handlers
	r.GET("/users", ListUsersHandler)
	r.POST("/users", CreateUserHandler)
	r.GET("/users/:id", GetUserHandler)
	r.GET("/ping", PingHandler)

	// Test /ping endpoint
	t.Run("GET /ping", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rec := httptest.NewRecorder()
		r.Handler().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "pong", resp["message"])
	})

	// Test /users (list) endpoint
	t.Run("GET /users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		r.Handler().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp struct {
			Data []User `json:"data"`
		}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Len(t, resp.Data, 2)
	})

	// Test /users (create) endpoint
	t.Run("POST /users", func(t *testing.T) {
		initialNextUserID := nextUserID // Save initial state
		initialUsers := make(map[int]User)
		for k, v := range users {
			initialUsers[k] = v
		}

		defer func() { // Restore initial state after test
			nextUserID = initialNextUserID
			users = initialUsers
		}()

		body := []byte(`{"name":"David","age":30}`)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.Handler().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		var resp User
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "David", resp.Name)
		assert.Equal(t, initialNextUserID, resp.ID)
	})

	// Test /users/:id (get) endpoint
	t.Run("GET /users/:id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec := httptest.NewRecorder()
		r.Handler().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp struct {
			Data User `json:"data"`
		}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "Alice", resp.Data.Name)
	})

	// Test 404 Not Found
	t.Run("GET /nonexistent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		rec := httptest.NewRecorder()
		r.Handler().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "Not Found", resp["error"])
	})
}

// Helper to set path parameters for unit tests
func (c *router.Context) SetParam(key, value string) {
	if c.params == nil {
		c.params = make(map[string]string)
	}
	c.params[key] = value
}
