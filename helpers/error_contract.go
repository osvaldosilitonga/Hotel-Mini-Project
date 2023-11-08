package helpers

import "net/http"

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  any    `json:"detail"`
}

var (
	// 400
	ErrBadRequest = APIError{
		Code:    http.StatusBadRequest,
		Message: "bad request",
	}
	ErrUnauthorized = APIError{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	}
	ErrNotFound = APIError{
		Code:    http.StatusNotFound,
		Message: "not found",
	}

	// 500
	ErrInternalServer = APIError{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
	}
)
