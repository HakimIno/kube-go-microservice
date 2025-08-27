package errors

// Error codes for different types of errors
const (
	// Authentication & Authorization
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired       = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid       = "TOKEN_INVALID"

	// Validation
	ErrCodeValidationFailed = "VALIDATION_FAILED"
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodeMissingRequired  = "MISSING_REQUIRED"
	ErrCodeInvalidFormat    = "INVALID_FORMAT"

	// Database
	ErrCodeDatabaseError       = "DATABASE_ERROR"
	ErrCodeRecordNotFound      = "RECORD_NOT_FOUND"
	ErrCodeDuplicateRecord     = "DUPLICATE_RECORD"
	ErrCodeConstraintViolation = "CONSTRAINT_VIOLATION"

	// Business Logic
	ErrCodeUserNotFound       = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists  = "USER_ALREADY_EXISTS"
	ErrCodeAccountDeactivated = "ACCOUNT_DEACTIVATED"
	ErrCodeInvalidOperation   = "INVALID_OPERATION"

	// External Services
	ErrCodeExternalServiceError = "EXTERNAL_SERVICE_ERROR"
	ErrCodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	ErrCodeTimeout              = "TIMEOUT"

	// System
	ErrCodeInternalError      = "INTERNAL_ERROR"
	ErrCodeConfigurationError = "CONFIGURATION_ERROR"
	ErrCodeRateLimitExceeded  = "RATE_LIMIT_EXCEEDED"
)

// HTTP status codes mapping
var ErrorCodeToHTTPStatus = map[string]int{
	ErrCodeUnauthorized:         401,
	ErrCodeForbidden:            403,
	ErrCodeInvalidCredentials:   401,
	ErrCodeTokenExpired:         401,
	ErrCodeTokenInvalid:         401,
	ErrCodeValidationFailed:     400,
	ErrCodeInvalidInput:         400,
	ErrCodeMissingRequired:      400,
	ErrCodeInvalidFormat:        400,
	ErrCodeDatabaseError:        500,
	ErrCodeRecordNotFound:       404,
	ErrCodeDuplicateRecord:      409,
	ErrCodeConstraintViolation:  400,
	ErrCodeUserNotFound:         404,
	ErrCodeUserAlreadyExists:    409,
	ErrCodeAccountDeactivated:   403,
	ErrCodeInvalidOperation:     400,
	ErrCodeExternalServiceError: 502,
	ErrCodeServiceUnavailable:   503,
	ErrCodeTimeout:              408,
	ErrCodeInternalError:        500,
	ErrCodeConfigurationError:   500,
	ErrCodeRateLimitExceeded:    429,
}
