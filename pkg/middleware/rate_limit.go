package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Simple cleanup goroutine to prevent memory leaks over time
	go func() {
		for {
			time.Sleep(window)
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	for ip, times := range rl.requests {
		// ⚡ Bolt Optimization: Slicing avoids allocating a new slice on every cleanup.
		// Since times are appended chronologically, we find the first valid one.
		i := 0
		for i < len(times) && !times[i].After(cutoff) {
			i++
		}

		if i == len(times) {
			// All times are expired
			delete(rl.requests, ip)
		} else {
			// Keep valid times
			rl.requests[ip] = times[i:]
		}
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	times, exists := rl.requests[ip]
	if !exists {
		// Initialize if not exists
		// Capacity set to a small number to reduce initial reallocations
		t := make([]time.Time, 0, 4)
		rl.requests[ip] = append(t, now)
		return true
	}

	// ⚡ Bolt Optimization: Slicing avoids allocating a new slice on every request.
	// Since times are appended chronologically, we just find the first valid one.
	cutoff := now.Add(-rl.window)
	i := 0
	for i < len(times) && !times[i].After(cutoff) {
		i++
	}

	times = times[i:]

	// Check if limit exceeded
	if len(times) >= rl.limit {
		rl.requests[ip] = times // keep valid times but deny
		return false
	}

	// Allow request and record time
	rl.requests[ip] = append(times, now)
	return true
}

// RateLimit returns a Gin middleware that limits the number of requests per IP address
// within a given time window.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := newRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "TOO_MANY_REQUESTS",
					"message": "Rate limit exceeded. Please try again later.",
				},
			})
			return
		}

		c.Next()
	}
}
