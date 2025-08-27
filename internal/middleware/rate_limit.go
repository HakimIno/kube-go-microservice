package middleware

import (
	"context"
	"sync"
	"time"

	"kube/pkg/errors"

	"github.com/cloudwego/hertz/pkg/app"
)

type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
	mutex    sync.RWMutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		mutex:    sync.RWMutex{},
	}
}

func (rl *RateLimiter) RateLimitMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		clientIP := string(c.ClientIP())
		now := time.Now()

		rl.mutex.Lock()
		defer rl.mutex.Unlock()

		// Clean up old requests
		if requests, exists := rl.requests[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < rl.window {
					validRequests = append(validRequests, reqTime)
				}
			}
			rl.requests[clientIP] = validRequests
		}

		// Check if limit exceeded
		if len(rl.requests[clientIP]) >= rl.limit {
			errors.SendError(c, errors.New(errors.ErrCodeRateLimitExceeded, "Rate limit exceeded", "Too many requests, please try again later"))
			c.Abort()
			return
		}

		// Add current request
		rl.requests[clientIP] = append(rl.requests[clientIP], now)
		c.Next(ctx)
	}
}

// AuthRateLimitMiddleware provides stricter rate limiting for authentication endpoints
func AuthRateLimitMiddleware(limit int, window time.Duration) app.HandlerFunc {
	limiter := NewRateLimiter(limit, window)
	return limiter.RateLimitMiddleware()
}
