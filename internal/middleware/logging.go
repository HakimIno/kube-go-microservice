package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// LogLevel represents different log levels
type LogLevel string

const (
	DEBUG   LogLevel = "DEBUG"
	INFO    LogLevel = "INFO"
	WARN    LogLevel = "WARN"
	ERROR   LogLevel = "ERROR"
	SUCCESS LogLevel = "SUCCESS"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level      LogLevel               `json:"level"`
	RequestID  string                 `json:"request_id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Path       string                 `json:"path,omitempty"`
	ClientIP   string                 `json:"client_ip,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	StatusCode int                    `json:"status_code,omitempty"`
	Duration   string                 `json:"duration,omitempty"`
	UserID     interface{}            `json:"user_id,omitempty"`
	Email      interface{}            `json:"email,omitempty"`
	Message    string                 `json:"message"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// Colors for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

// getColorByLevel returns the appropriate color for each log level
func getColorByLevel(level LogLevel) string {
	switch level {
	case DEBUG:
		return colorCyan
	case INFO:
		return colorBlue
	case WARN:
		return colorYellow
	case ERROR:
		return colorRed
	case SUCCESS:
		return colorGreen
	default:
		return colorWhite
	}
}

// formatDuration formats duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fμs", float64(d.Nanoseconds())/1000.0)
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Milliseconds()))
	}
	return fmt.Sprintf("%.2fs", float64(d.Seconds()))
}

// logStructured logs a structured log entry
func logStructured(entry LogEntry) {
	// Check if we're in a terminal environment
	if os.Getenv("TERM") != "" && os.Getenv("NO_COLOR") == "" {
		// Colored terminal output
		color := getColorByLevel(entry.Level)
		levelStr := fmt.Sprintf("%s%s%s", color, entry.Level, colorReset)

		// Format the log message with colors
		message := fmt.Sprintf("%s[%s]%s %s %s",
			colorBold,
			entry.Timestamp,
			colorReset,
			levelStr,
			entry.Message,
		)

		// Add request details if available
		if entry.RequestID != "" {
			message += fmt.Sprintf(" %s[%s]%s", colorPurple, entry.RequestID, colorReset)
		}

		// Add method and path if available
		if entry.Method != "" && entry.Path != "" {
			message += fmt.Sprintf(" %s%s %s%s", colorCyan, entry.Method, entry.Path, colorReset)
		}

		// Add status code and duration if available
		if entry.StatusCode > 0 {
			statusColor := colorGreen
			if entry.StatusCode >= 400 {
				statusColor = colorRed
			} else if entry.StatusCode >= 300 {
				statusColor = colorYellow
			}
			message += fmt.Sprintf(" %s[%d]%s", statusColor, entry.StatusCode, colorReset)
		}

		if entry.Duration != "" {
			message += fmt.Sprintf(" %s(%s)%s", colorYellow, entry.Duration, colorReset)
		}

		// Add user info if available
		if entry.UserID != nil {
			message += fmt.Sprintf(" %sUser:%v%s", colorGreen, entry.UserID, colorReset)
		}

		fmt.Println(message)
	} else {
		// JSON output for non-terminal environments
		jsonData, _ := json.Marshal(entry)
		fmt.Println(string(jsonData))
	}
}

// LoggingMiddleware provides beautiful structured logging
func LoggingMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Request.URI().Path())
		method := string(c.Method())
		clientIP := string(c.ClientIP())
		userAgent := string(c.GetHeader("User-Agent"))
		requestID := c.GetString("request_id")

		// Log request start
		logStructured(LogEntry{
			Timestamp: time.Now().Format("15:04:05.000"),
			Level:     INFO,
			RequestID: requestID,
			Method:    method,
			Path:      path,
			ClientIP:  clientIP,
			UserAgent: userAgent,
			Message:   "Request started",
		})

		c.Next(ctx)

		duration := time.Since(start)
		statusCode := c.Response.StatusCode()

		// Get user info if available
		userID, _ := c.Get("user_id")
		email, _ := c.Get("email")

		// Determine log level based on status code
		level := INFO
		if statusCode >= 500 {
			level = ERROR
		} else if statusCode >= 400 {
			level = WARN
		} else if statusCode >= 200 && statusCode < 300 {
			level = SUCCESS
		}

		// Log response
		logStructured(LogEntry{
			Timestamp:  time.Now().Format("15:04:05.000"),
			Level:      level,
			RequestID:  requestID,
			Method:     method,
			Path:       path,
			StatusCode: statusCode,
			Duration:   formatDuration(duration),
			UserID:     userID,
			Email:      email,
			Message:    "Request completed",
		})

		// Security logging for failed authentication attempts
		if (path == "/api/v1/auth/login" || path == "/api/v1/users/register") && statusCode >= 400 {
			logStructured(LogEntry{
				Timestamp:  time.Now().Format("15:04:05.000"),
				Level:      WARN,
				RequestID:  requestID,
				Method:     method,
				Path:       path,
				ClientIP:   clientIP,
				StatusCode: statusCode,
				Message:    "Security Alert: Failed authentication attempt",
				Extra: map[string]interface{}{
					"security_event": "failed_auth",
					"ip_address":     clientIP,
				},
			})
		}
	}
}

// SecurityLoggingMiddleware logs security-related events with beautiful formatting
func SecurityLoggingMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())
		method := string(c.Method())
		clientIP := string(c.ClientIP())
		requestID := c.GetString("request_id")

		// Log authentication attempts
		if path == "/api/v1/auth/login" {
			logStructured(LogEntry{
				Timestamp: time.Now().Format("15:04:05.000"),
				Level:     INFO,
				RequestID: requestID,
				Method:    method,
				Path:      path,
				ClientIP:  clientIP,
				Message:   "Authentication attempt",
				Extra: map[string]interface{}{
					"security_event": "auth_attempt",
				},
			})
		}

		// Log registration attempts
		if path == "/api/v1/users/register" {
			logStructured(LogEntry{
				Timestamp: time.Now().Format("15:04:05.000"),
				Level:     INFO,
				RequestID: requestID,
				Method:    method,
				Path:      path,
				ClientIP:  clientIP,
				Message:   "Registration attempt",
				Extra: map[string]interface{}{
					"security_event": "registration_attempt",
				},
			})
		}

		// Log password change attempts
		if path == "/api/v1/auth/change-password" {
			userID, hasUser := c.Get("user_id")
			if hasUser {
				logStructured(LogEntry{
					Timestamp: time.Now().Format("15:04:05.000"),
					Level:     INFO,
					RequestID: requestID,
					Method:    method,
					Path:      path,
					ClientIP:  clientIP,
					UserID:    userID,
					Message:   "Password change attempt",
					Extra: map[string]interface{}{
						"security_event": "password_change_attempt",
					},
				})
			}
		}

		c.Next(ctx)

		// Log failed attempts
		statusCode := c.Response.StatusCode()
		if statusCode >= 400 && (path == "/api/v1/auth/login" || path == "/api/v1/users/register") {
			logStructured(LogEntry{
				Timestamp:  time.Now().Format("15:04:05.000"),
				Level:      WARN,
				RequestID:  requestID,
				Method:     method,
				Path:       path,
				ClientIP:   clientIP,
				StatusCode: statusCode,
				Message:    "Failed authentication/registration",
				Extra: map[string]interface{}{
					"security_event": "failed_auth_registration",
					"ip_address":     clientIP,
				},
			})
		}
	}
}

// DatabaseLoggingMiddleware logs database operations
func DatabaseLoggingMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		requestID := c.GetString("request_id")

		logStructured(LogEntry{
			Timestamp: time.Now().Format("15:04:05.000"),
			Level:     DEBUG,
			RequestID: requestID,
			Message:   "Database operation",
			Extra: map[string]interface{}{
				"operation": "db_query",
			},
		})

		c.Next(ctx)
	}
}
