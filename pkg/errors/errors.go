package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Error codes
const (
	CodeNotFound     = "NOT_FOUND"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeValidation   = "VALIDATION_ERROR"
	CodeInternal     = "INTERNAL_ERROR"
	CodeConflict     = "CONFLICT"
	CodeBadRequest   = "BAD_REQUEST"
)

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    CodeUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    CodeForbidden,
		Message: message,
	}
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    CodeValidation,
		Message: message,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Code:    CodeInternal,
		Message: message,
		Err:     err,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    CodeConflict,
		Message: message,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    CodeBadRequest,
		Message: message,
	}
}

