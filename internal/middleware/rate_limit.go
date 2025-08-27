package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) RateLimitMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		clientIP := string(c.ClientIP())
		now := time.Now()

		if requests, exists := rl.requests[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < rl.window {
					validRequests = append(validRequests, reqTime)
				}
			}
			rl.requests[clientIP] = validRequests
		}

		if len(rl.requests[clientIP]) >= rl.limit {
			c.JSON(429, utils.H{
				"error":       "Rate limit exceeded",
				"retry_after": rl.window.Seconds(),
			})
			c.Abort()
			return
		}

		rl.requests[clientIP] = append(rl.requests[clientIP], now)
		c.Next(ctx)
	}
}
