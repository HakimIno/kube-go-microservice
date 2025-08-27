package errors

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

type ErrorResponse struct {
	Success   bool      `json:"success"`
	Error     *AppError `json:"error"`
	Timestamp string    `json:"timestamp"`
	Path      string    `json:"path,omitempty"`
	Method    string    `json:"method,omitempty"`
}

func SendError(c *app.RequestContext, err error) {
	var appErr *AppError

	if IsAppError(err) {
		appErr = GetAppError(err)
	} else {
		appErr = Wrap(err, ErrCodeInternalError, "Internal server error", err.Error())
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		appErr.RequestID = requestID
	}

	response := ErrorResponse{
		Success:   false,
		Error:     appErr,
		Timestamp: appErr.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		Path:      string(c.Request.URI().Path()),
		Method:    string(c.Request.Method()),
	}

	c.JSON(appErr.HTTPStatus, response)
}

func SendValidationError(c *app.RequestContext, details string) {
	err := New(ErrCodeValidationFailed, "Validation failed", details)
	SendError(c, err)
}

func SendNotFoundError(c *app.RequestContext, resource string) {
	err := New(ErrCodeRecordNotFound, "Resource not found", resource+" not found")
	SendError(c, err)
}

func SendUnauthorizedError(c *app.RequestContext, details string) {
	err := New(ErrCodeUnauthorized, "Unauthorized", details)
	SendError(c, err)
}

func SendForbiddenError(c *app.RequestContext, details string) {
	err := New(ErrCodeForbidden, "Forbidden", details)
	SendError(c, err)
}

func SendDatabaseError(c *app.RequestContext, err error) {
	appErr := Wrap(err, ErrCodeDatabaseError, "Database operation failed", err.Error())
	SendError(c, appErr)
}

func SendRateLimitError(c *app.RequestContext) {
	err := New(ErrCodeRateLimitExceeded, "Rate limit exceeded", "Too many requests")
	SendError(c, err)
}

func SendSuccess(c *app.RequestContext, statusCode int, data interface{}, message string) {
	response := map[string]interface{}{
		"success":   true,
		"message":   message,
		"data":      data,
		"timestamp": time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		response["request_id"] = requestID
	}

	c.JSON(statusCode, response)
}
