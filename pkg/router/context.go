package router

import (
    "encoding/json"
    "net/http"
    "sync"
)

type Context struct {
    req        *http.Request
    res        http.ResponseWriter
    params     map[string]string
    store      map[string]any
    mu         sync.RWMutex
    statusCode int
    written    bool
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
    return &Context{
        req:        r,
        res:        w,
        params:     make(map[string]string),
        store:      make(map[string]any),
        statusCode: http.StatusOK,
    }
}

func (c *Context) Request() *http.Request {
    return c.req
}

func (c *Context) ResponseWriter() http.ResponseWriter {
    return c.res
}

// --------- ROUTE PARAMS ---------

func (c *Context) Param(name string) string {
    return c.params[name]
}

// --------- QUERY / FORM ---------

func (c *Context) Query(name string) string {
    return c.req.URL.Query().Get(name)
}

func (c *Context) FormValue(name string) string {
    return c.req.FormValue(name)
}

// --------- JSON BINDING ---------

func (c *Context) BindJSON(dest any) error {
    return json.NewDecoder(c.req.Body).Decode(dest)
}

// --------- RESPONSE HELPERS ---------

func (c *Context) Status(code int) {
    if c.written {
        return
    }
    c.statusCode = code
    c.res.WriteHeader(code)
    c.written = true
}

func (c *Context) JSON(status int, value any) {
    c.res.Header().Set("Content-Type", "application/json")
    c.Status(status)
    if err := json.NewEncoder(c.res).Encode(value); err != nil {
        http.Error(c.res, err.Error(), http.StatusInternalServerError)
    }
}

// --------- CONTEXT STORE ---------

func (c *Context) Set(key string, value any) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.store[key] = value
}

func (c *Context) Get(key string) (any, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.store[key]
    return v, ok
}

func (c *Context) GetString(key string) (string, bool) {
    v, ok := c.Get(key)
    if !ok {
        return "", false
    }
    s, ok := v.(string)
    return s, ok
}

func (c *Context) GetInt(key string) (int, bool) {
    v, ok := c.Get(key)
    if !ok {
        return 0, false
    }
    i, ok := v.(int)
    return i, ok
}
