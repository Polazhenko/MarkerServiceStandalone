package models

// Error is a generic error response struct
type Error struct {
	Message string `json:"error"`
}

// ErrorResponse wraps error details for API responses
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"`
}

// NewError creates a new error response
func NewError(message string) *Error {
	return &Error{Message: message}
}

// NewErrorResponse creates a new error response with code
func NewErrorResponse(message string, code int) *ErrorResponse {
	return &ErrorResponse{
		Error: message,
		Code:  code,
	}
}
