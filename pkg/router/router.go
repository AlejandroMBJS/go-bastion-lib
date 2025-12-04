// Package router provides a fast, expressive HTTP router with middleware support.
package router

import (
	"net/http"
	"strings"
	"sync"
)

// Handler is the function signature for HTTP handlers.
type Handler func(ctx *Context)

type routeMatch struct {
    handler Handler
    params  map[string]string
}


// Middleware is a function that wraps a Handler.
type Middleware func(Handler) Handler

// Router manages routes and middleware.
type Router struct {
	prefix           string
	parent           *Router
	tree             *node
	middlewares      []Middleware
	mu               sync.RWMutex
	notFound         Handler
	methodNotAllowed Handler
}

// node represents a node in the radix tree.
type node struct {
	path      string
	isParam   bool
	paramName string
	children  []*node
	handlers  map[string]Handler
}

// New creates a new Router instance.
func New() *Router {
	r := &Router{
		tree:        &node{},
		middlewares: []Middleware{},
		notFound: func(ctx *Context) {
			ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "Not Found",
			})
		},
		methodNotAllowed: func(ctx *Context) {
			ctx.JSON(http.StatusMethodNotAllowed, map[string]string{
				"error": "Method Not Allowed",
			})
		},
	}
	return r
}

// Group creates a new router group with the given prefix.
func (r *Router) Group(prefix string) *Router {
	return &Router{
		prefix:           r.prefix + prefix,
		parent:           r,
		tree:             r.tree,
		middlewares:      r.middlewares,
		notFound:         r.notFound,
		methodNotAllowed: r.methodNotAllowed,
	}
}

// Use registers middleware that will be applied to all routes in this group.
func (r *Router) Use(mw ...Middleware) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewares = append(r.middlewares, mw...)
}

// GET registers a GET route.
func (r *Router) GET(path string, h Handler) {
	r.addRoute(http.MethodGet, path, h)
}

// POST registers a POST route.
func (r *Router) POST(path string, h Handler) {
	r.addRoute(http.MethodPost, path, h)
}

// PUT registers a PUT route.
func (r *Router) PUT(path string, h Handler) {
	r.addRoute(http.MethodPut, path, h)
}

// DELETE registers a DELETE route.
func (r *Router) DELETE(path string, h Handler) {
	r.addRoute(http.MethodDelete, path, h)
}

// PATCH registers a PATCH route.
func (r *Router) PATCH(path string, h Handler) {
	r.addRoute(http.MethodPatch, path, h)
}

// Handler returns the HTTP handler for the router.
func (r *Router) Handler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        ctx := NewContext(w, req)

        match := r.findRoute(req.Method, req.URL.Path)
        if match.handler == nil {
            // Check if it's a method-not-allowed or a true 404
            if r.isPathRegistered(req.URL.Path) {
                r.methodNotAllowed(ctx)
            } else {
                r.notFound(ctx)
            }
            return
        }

        // inject path params into context
        for k, v := range match.params {
            ctx.params[k] = v
        }

        match.handler(ctx)
    })
}


// addRoute adds a route with the given method and path.
func (r *Router) addRoute(method, path string, h Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fullPath := r.prefix + path
	if !strings.HasPrefix(fullPath, "/") {
		fullPath = "/" + fullPath
	}

	// Apply middlewares to the handler
	wrappedHandler := h
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		wrappedHandler = r.middlewares[i](wrappedHandler)
	}

	r.tree.insert(method, fullPath, wrappedHandler)
}

// findRoute finds a handler for the given method and path.
func (r *Router) findRoute(method, path string) routeMatch {
    r.mu.RLock()
    defer r.mu.RUnlock()

    if r.tree == nil {
        return routeMatch{}
    }

    return r.tree.find(method, path)
}


// isPathRegistered checks if a path is registered (for any method).
func (r *Router) isPathRegistered(path string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut,
		http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		if r.tree.find(method, path) != nil {
			return true
		}
	}
	return false
}

// node methods
func (n *node) insert(method, path string, handler Handler) {
	if n.children == nil {
		n.children = []*node{}
	}

	// Split path into segments
	segments := strings.Split(strings.Trim(path, "/"), "/")
	n.insertRecursive(method, segments, handler)
}

func (n *node) insertRecursive(method string, segments []string, handler Handler) {
	if len(segments) == 0 {
		if n.handlers == nil {
			n.handlers = make(map[string]Handler)
		}
		n.handlers[method] = handler
		return
	}

	segment := segments[0]
	var child *node

	// Check for existing child
	for _, c := range n.children {
		if c.path == segment {
			child = c
			break
		}
	}

	// Create new child if not found
	if child == nil {
		child = &node{
			path:    segment,
			isParam: strings.HasPrefix(segment, ":"),
		}
		if child.isParam {
			child.paramName = segment[1:]
		}
		n.children = append(n.children, child)
	}

	child.insertRecursive(method, segments[1:], handler)
}

// find returns the handler for a method+path and any path parameters.
func (n *node) find(method, path string) routeMatch {
    segments := strings.Split(strings.Trim(path, "/"), "/")
    params := make(map[string]string)

    h := n.findRecursive(method, segments, params)
    if h == nil {
        return routeMatch{}
    }

    return routeMatch{
        handler: h,
        params:  params,
    }
}

func (n *node) findRecursive(method string, segments []string, params map[string]string) Handler {
    if len(segments) == 0 {
        if n.handlers == nil {
            return nil
        }
        return n.handlers[method]
    }

    segment := segments[0]

    // 1) exact match first
    for _, child := range n.children {
        if !child.isParam && child.path == segment {
            if h := child.findRecursive(method, segments[1:], params); h != nil {
                return h
            }
        }
    }

    // 2) then param match
    for _, child := range n.children {
        if child.isParam {
            params[child.paramName] = segment
            if h := child.findRecursive(method, segments[1:], params); h != nil {
                return h
            }
        }
    }

    return nil
}

}
