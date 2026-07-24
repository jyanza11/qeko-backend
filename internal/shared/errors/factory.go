package errors

import "net/http"

func Validation(details any) *AppError {

	return &AppError{
		Code:       CodeValidation,
		Message:    "Validation error",
		StatusCode: http.StatusBadRequest,
		Details:    details,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       CodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Code:       CodeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:       CodeNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Code:       CodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func Internal(message string) *AppError {
	return &AppError{
		Code:       CodeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}
