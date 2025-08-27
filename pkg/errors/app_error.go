package errors

import (
	"fmt"
	"time"
)

type AppError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	HTTPStatus int                    `json:"-"`
	Timestamp  time.Time              `json:"timestamp"`
	RequestID  string                 `json:"request_id,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Original   error                  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Original != nil {
		return fmt.Sprintf("[%s] %s: %s (original: %v)", e.Code, e.Message, e.Details, e.Original)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
}

func (e *AppError) Unwrap() error {
	return e.Original
}

func New(code, message, details string) *AppError {
	httpStatus := ErrorCodeToHTTPStatus[code]
	if httpStatus == 0 {
		httpStatus = 500
	}

	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: httpStatus,
		Timestamp:  time.Now(),
	}
}

func NewWithRequestID(code, message, details, requestID string) *AppError {
	err := New(code, message, details)
	err.RequestID = requestID
	return err
}

func Wrap(err error, code, message, details string) *AppError {
	appErr := New(code, message, details)
	appErr.Original = err
	return appErr
}

func WrapWithRequestID(err error, code, message, details, requestID string) *AppError {
	appErr := Wrap(err, code, message, details)
	appErr.RequestID = requestID
	return appErr
}

func (e *AppError) AddMetadata(key string, value interface{}) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}
