package server

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		requestID := string(c.GetHeader("X-Request-ID"))
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next(ctx)
	}
}

// LoggingMiddleware logs request details
func LoggingMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		requestID := c.GetString("request_id")

		hlog.Infof("[%s] %s %s - Started", requestID, c.Method(), c.Request.URI().String())

		c.Next(ctx)

		duration := time.Since(start)
		status := c.Response.StatusCode()

		hlog.Infof("[%s] %s %s - %d - %v", requestID, c.Method(), c.Request.URI().String(), status, duration)
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(consts.StatusInternalServerError, utils.H{
					"error":   "Internal server error",
					"message": "Something went wrong",
				})
			}
		}()
		c.Next(ctx)
	}
}

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requests int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Every(duration/time.Duration(requests)), requests),
	}
}

// RateLimitMiddleware applies rate limiting
func (rl *RateLimiter) RateLimitMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !rl.limiter.Allow() {
			c.JSON(consts.StatusTooManyRequests, utils.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			return
		}
		c.Next(ctx)
	}
}
