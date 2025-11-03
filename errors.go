package prometheus

import (
	"net/http"
)

// ErrorType classifies HTTP errors for better observability
type ErrorType string

const (
	ErrorTypeClientError ErrorType = "client_error" // 4xx errors
	ErrorTypeServerError ErrorType = "server_error" // 5xx errors
	ErrorTypeTimeout     ErrorType = "timeout"      // 408, 504
	ErrorTypeNoWorkers   ErrorType = "no_workers"   // Worker pool exhausted
)

// classifyError determines the error type based on status code and headers
func classifyError(statusCode int, headers http.Header) ErrorType {
	// Check for no workers condition first
	if headers.Get(noWorkers) == trueStr {
		return ErrorTypeNoWorkers
	}

	// Classify by status code
	switch {
	case statusCode == 408 || statusCode == 504:
		return ErrorTypeTimeout
	case statusCode >= 400 && statusCode < 500:
		return ErrorTypeClientError
	case statusCode >= 500:
		return ErrorTypeServerError
	default:
		return ErrorTypeServerError
	}
}

// isErrorStatus checks if a status code represents an error
func isErrorStatus(statusCode int) bool {
	return statusCode >= 400
}
