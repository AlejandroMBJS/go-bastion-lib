package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

// RateLimit creates a rate limiting middleware.
func RateLimit(requests int, window time.Duration) router.Middleware {
	limiter := newSlidingWindowLimiter(requests, window)

	return func(next router.Handler) router.Handler {
		return func(ctx *router.Context) {
			ip := getClientIP(ctx.Request())

			if !limiter.Allow(ip) {
				ctx.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "too_many_requests",
				})
				return
			}

			next(ctx)
		}
	}
}

// slidingWindowLimiter implements a sliding window rate limiter.
type slidingWindowLimiter struct {
	requests int
	window   time.Duration
	visits   map[string][]time.Time
	mu       sync.RWMutex
}

// newSlidingWindowLimiter creates a new sliding window rate limiter.
func newSlidingWindowLimiter(requests int, window time.Duration) *slidingWindowLimiter {
	limiter := &slidingWindowLimiter{
		requests: requests,
		window:   window,
		visits:   make(map[string][]time.Time),
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// Allow checks if the given key is allowed to make a request.
func (l *slidingWindowLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	// Get visits for this key
	visits := l.visits[key]

	// Remove old visits
	i := 0
	for i < len(visits) && visits[i].Before(cutoff) {
		i++
	}
	visits = visits[i:]

	// Check if we've exceeded the limit
	if len(visits) >= l.requests {
		return false
	}

	// Add new visit
	visits = append(visits, now)
	l.visits[key] = visits

	return true
}

// cleanup periodically removes old entries.
func (l *slidingWindowLimiter) cleanup() {
	ticker := time.NewTicker(l.window)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		cutoff := time.Now().Add(-l.window)

		for key, visits := range l.visits {
			// Remove old visits
			i := 0
			for i < len(visits) && visits[i].Before(cutoff) {
				i++
			}

			if i == len(visits) {
				// All visits are old, remove the key
				delete(l.visits, key)
			} else {
				// Keep only recent visits
				l.visits[key] = visits[i:]
			}
		}
		l.mu.Unlock()
	}
}

// getClientIP extracts the client IP from the request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxy)
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}

	// Check X-Real-IP header
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// Fall back to remote address
	return r.RemoteAddr
}
