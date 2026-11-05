package domain

import "fmt"

type ErrorCode string

const (
	CodeNotFound         ErrorCode = "NOT_FOUND"
	CodeInvalidArgument  ErrorCode = "INVALID_ARGUMENT"
	CodePermissionDenied ErrorCode = "PERMISSION_DENIED"
	CodeAlreadyExists    ErrorCode = "ALREADY_EXISTS"
	CodeUnauthenticated  ErrorCode = "UNAUTHENTICATED"
	CodeInternal         ErrorCode = "INTERNAL"
)

type AppError struct {
	Code    ErrorCode
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(code ErrorCode, format string, args ...any) error {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
