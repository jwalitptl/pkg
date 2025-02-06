package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common errors
func NotFound(resource string, err error) *AppError {
	return NewAppError(
		http.StatusNotFound,
		fmt.Sprintf("%s not found", resource),
		err,
	)
}

func BadRequest(message string, err error) *AppError {
	return NewAppError(
		http.StatusBadRequest,
		message,
		err,
	)
}

func Internal(err error) *AppError {
	return NewAppError(
		http.StatusInternalServerError,
		"internal server error",
		err,
	)
}

func Unauthorized(err error) *AppError {
	return NewAppError(
		http.StatusUnauthorized,
		"unauthorized",
		err,
	)
}
