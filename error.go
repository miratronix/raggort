package raggort

import (
	"net/http"
)

// NewError creates a new error with the specified status code and message
func NewError(statusCode int, name string, message string) *HTTPResponse {
	return NewHTTPResponse().Status(statusCode).Body(map[string]string{
		"name":    name,
		"message": message,
	})
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string) *HTTPResponse {
	return NewError(http.StatusBadRequest, "BadRequestError", message)
}

// NewUnauthorizedError creates a new 401 unauthorized error
func NewUnauthorizedError() *HTTPResponse {
	return NewError(http.StatusUnauthorized, "UnauthorizedError", "Unauthorized")
}

// NewForbiddenError creates a new 403 forbidden error
func NewForbiddenError() *HTTPResponse {
	return NewError(http.StatusForbidden, "ForbiddenError", "Forbidden")
}

// NewNotFoundError creates a new 404 error
func NewNotFoundError() *HTTPResponse {
	return NewError(http.StatusNotFound, "NotFoundError", "The request resource could not be found")
}

// NewConflictError creates a new 409 error
func NewConflictError(message string) *HTTPResponse {
	return NewError(http.StatusConflict, "ConflictError", message)
}

// NewUnsupportedMediaTypeError creates a new 415 error
func NewUnsupportedMediaTypeError() *HTTPResponse {
	return NewError(http.StatusUnsupportedMediaType, "UnsupportedMediaTypeError", "Unsupported Media Type")
}

// NewTooManyRequestsError creates a new 429 error
func NewTooManyRequestsError(message string) *HTTPResponse {
	return NewError(http.StatusTooManyRequests, "TooManyRequestsError", message)
}

// NewInternalServerError creates a new 500 error
func NewInternalServerError(message string) *HTTPResponse {
	return NewError(http.StatusInternalServerError, "InternalServerError", message)
}
