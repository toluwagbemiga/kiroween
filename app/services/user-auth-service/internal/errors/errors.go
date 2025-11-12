package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode represents an error code
type ErrorCode string

const (
	ErrCodeInvalidCredentials  ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeUserNotFound        ErrorCode = "USER_NOT_FOUND"
	ErrCodeEmailAlreadyExists  ErrorCode = "EMAIL_ALREADY_EXISTS"
	ErrCodeInvalidToken        ErrorCode = "INVALID_TOKEN"
	ErrCodeExpiredToken        ErrorCode = "EXPIRED_TOKEN"
	ErrCodeRevokedToken        ErrorCode = "REVOKED_TOKEN"
	ErrCodeAccountLocked       ErrorCode = "ACCOUNT_LOCKED"
	ErrCodePermissionDenied    ErrorCode = "PERMISSION_DENIED"
	ErrCodeInvalidInput        ErrorCode = "INVALID_INPUT"
	ErrCodeRoleNotFound        ErrorCode = "ROLE_NOT_FOUND"
	ErrCodeSystemRoleProtected ErrorCode = "SYSTEM_ROLE_PROTECTED"
	ErrCodeInternal            ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidEmail        ErrorCode = "INVALID_EMAIL"
	ErrCodeWeakPassword        ErrorCode = "WEAK_PASSWORD"
	ErrCodeInvalidResetToken   ErrorCode = "INVALID_RESET_TOKEN"
)

// ServiceError represents a service-level error
type ServiceError struct {
	Code    ErrorCode
	Message string
	Details map[string]interface{}
	Cause   error
}

// Error implements the error interface
func (e *ServiceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// New creates a new ServiceError
func New(code ErrorCode, message string) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// Wrap wraps an error with a ServiceError
func Wrap(code ErrorCode, message string, cause error) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *ServiceError) WithDetails(key string, value interface{}) *ServiceError {
	e.Details[key] = value
	return e
}

// MapToGRPCError maps a ServiceError to a gRPC error
func MapToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	serviceErr, ok := err.(*ServiceError)
	if !ok {
		return status.Error(codes.Internal, "internal server error")
	}

	switch serviceErr.Code {
	case ErrCodeInvalidCredentials, ErrCodeInvalidToken, ErrCodeExpiredToken, ErrCodeRevokedToken:
		return status.Error(codes.Unauthenticated, serviceErr.Message)
	case ErrCodePermissionDenied, ErrCodeAccountLocked:
		return status.Error(codes.PermissionDenied, serviceErr.Message)
	case ErrCodeUserNotFound, ErrCodeRoleNotFound:
		return status.Error(codes.NotFound, serviceErr.Message)
	case ErrCodeEmailAlreadyExists:
		return status.Error(codes.AlreadyExists, serviceErr.Message)
	case ErrCodeInvalidInput, ErrCodeInvalidEmail, ErrCodeWeakPassword, ErrCodeInvalidResetToken:
		return status.Error(codes.InvalidArgument, serviceErr.Message)
	case ErrCodeSystemRoleProtected:
		return status.Error(codes.FailedPrecondition, serviceErr.Message)
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
