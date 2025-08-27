package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func LoggingMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Request.URI().Path())
		method := string(c.Method())
		clientIP := string(c.ClientIP())
		userAgent := string(c.GetHeader("User-Agent"))

		// Log request
		hlog.Infof("Request: %s %s from %s (User-Agent: %s)",
			method,
			path,
			clientIP,
			userAgent,
		)

		c.Next(ctx)

		duration := time.Since(start)
		statusCode := c.Response.StatusCode()

		// Get user info if available
		userID, hasUser := c.Get("user_id")
		email, hasEmail := c.Get("email")

		if hasUser && hasEmail {
			hlog.Infof("Response: %s %s - %d - %v - User: %v (%s)",
				method,
				path,
				statusCode,
				duration,
				userID,
				email,
			)
		} else {
			hlog.Infof("Response: %s %s - %d - %v - Anonymous",
				method,
				path,
				statusCode,
				duration,
			)
		}

		// Security logging for failed authentication attempts
		if (path == "/api/v1/users/login" || path == "/api/v1/users/register") && statusCode >= 400 {
			hlog.Warnf("Security Alert: Failed authentication attempt - %s %s from %s - Status: %d",
				method,
				path,
				clientIP,
				statusCode,
			)
		}
	}
}

// SecurityLoggingMiddleware logs security-related events
func SecurityLoggingMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())
		method := string(c.Method())
		clientIP := string(c.ClientIP())

		// Log authentication attempts
		if path == "/api/v1/users/login" {
			hlog.Infof("Authentication attempt: %s from %s", method, clientIP)
		}

		// Log registration attempts
		if path == "/api/v1/users/register" {
			hlog.Infof("Registration attempt: %s from %s", method, clientIP)
		}

		// Log password change attempts
		if path == "/api/v1/users/change-password" {
			userID, hasUser := c.Get("user_id")
			if hasUser {
				hlog.Infof("Password change attempt: User %v from %s", userID, clientIP)
			}
		}

		c.Next(ctx)

		// Log failed attempts
		statusCode := c.Response.StatusCode()
		if statusCode >= 400 && (path == "/api/v1/users/login" || path == "/api/v1/users/register") {
			hlog.Warnf("Failed authentication/registration: %s %s from %s - Status: %d",
				method,
				path,
				clientIP,
				statusCode,
			)
		}
	}
}
