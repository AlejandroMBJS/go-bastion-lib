package router

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Context represents the request context.
type Context struct {
	req        *http.Request
	res        http.ResponseWriter
	params     map[string]string
	store      map[string]any
	mu         sync.RWMutex
	statusCode int
	written    bool
}

// NewContext creates a new Context instance.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		req:        r,
		res:        w,
		params:     make(map[string]string),
		store:      make(map[string]any),
		statusCode: http.StatusOK,
	}
}

// Request returns the underlying http.Request.
func (c *Context) Request() *http.Request {
	return c.req
}

// ResponseWriter returns the underlying http.ResponseWriter.
func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.res
}

// Param returns the value of a path parameter.
func (c *Context) Param(name string) string {
	return c.params[name]
}

// Query returns the value of a query parameter.
func (c *Context) Query(name string) string {
	return c.req.URL.Query().Get(name)
}

// FormValue returns the value of a form parameter.
func (c *Context) FormValue(name string) string {
	return c.req.FormValue(name)
}

// BindJSON decodes the JSON request body into the given value.
func (c *Context) BindJSON(dest any) error {
	return json.NewDecoder(c.req.Body).Decode(dest)
}

// Status sets the HTTP status code.
func (c *Context) Status(code int) {
	c.statusCode = code
	c.res.WriteHeader(code)
	c.written = true
}

// JSON sends a JSON response with the given status code.
func (c *Context) JSON(status int, value any) {
	c.res.Header().Set("Content-Type", "application/json")
	c.Status(status)

	if err := json.NewEncoder(c.res).Encode(value); err != nil {
		// If encoding fails, we can't send JSON error since we already wrote status
		http.Error(c.res, err.Error(), http.StatusInternalServerError)
	}
}

// Set stores a value in the context.
func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

// Get retrieves a value from the context.
func (c *Context) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.store[key]
	return value, ok
}

// GetString retrieves a string value from the context.
func (c *Context) GetString(key string) (string, bool) {
	if value, ok := c.Get(key); ok {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetInt retrieves an int value from the context.
func (c *Context) GetInt(key string) (int, bool) {
	if value, ok := c.Get(key); ok {
		switch v := value.(type) {
		case int:
			return v, true
		case float64:
			return int(v), true
		}
	}
	return 0, false
}

// StatusCode returns the HTTP status code set for the context.
func (c *Context) StatusCode() int {
	return c.statusCode
}
