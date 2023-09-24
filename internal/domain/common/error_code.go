package common

import "net/http"

// ErrorCode represents a custom error code structure, containing
// a name and a status code that can be used for HTTP responses.
type ErrorCode struct {
	Name       string
	StatusCode int
}

// ErrorCodeInternalProcess represents a general internal process error.
// This error will be associated with an HTTP Internal Server Error (500) status.
var ErrorCodeInternalProcess = ErrorCode{
	Name:       "INTERNAL_PROCESS",
	StatusCode: http.StatusInternalServerError,
}

// ErrorCodeAuthPermissionDenied represents an authentication error where permission is denied.
// This error will be associated with an HTTP Forbidden (403) status.
var ErrorCodeAuthPermissionDenied = ErrorCode{
	Name:       "AUTH_PERMISSION_DENIED",
	StatusCode: http.StatusForbidden,
}

// ErrorCodeAuthNotAuthenticated represents an authentication error where the user is not authenticated.
// This error will be associated with an HTTP Unauthorized (401) status.
var ErrorCodeAuthNotAuthenticated = ErrorCode{
	Name:       "AUTH_NOT_AUTHENTICATED",
	StatusCode: http.StatusUnauthorized,
}

// ErrorCodeResourceNotFound represents an error where the requested resource is not found.
// This error will be associated with an HTTP Not Found (404) status.
var ErrorCodeResourceNotFound = ErrorCode{
	Name:       "RESOURCE_NOT_FOUND",
	StatusCode: http.StatusNotFound,
}

// ErrorCodeParameterInvalid represents an error where the provided parameter is invalid.
// This error will be associated with an HTTP Bad Request (400) status.
var ErrorCodeParameterInvalid = ErrorCode{
	Name:       "PARAMETER_INVALID",
	StatusCode: http.StatusBadRequest,
}

// ErrorCodeRemoteProcess represents an error in the remote process.
// This error will be associated with an HTTP Bad Gateway (502) status.
var ErrorCodeRemoteProcess = ErrorCode{
	Name:       "REMOTE_PROCESS_ERROR",
	StatusCode: http.StatusBadGateway,
}
