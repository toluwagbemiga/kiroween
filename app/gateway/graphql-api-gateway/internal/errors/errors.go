package errors

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ConvertGRPCError converts a gRPC error to a user-friendly GraphQL error
func ConvertGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		// Not a gRPC error - return generic error
		return &gqlerror.Error{
			Message: "An unexpected error occurred",
			Extensions: map[string]interface{}{
				"code": "INTERNAL_ERROR",
			},
		}
	}

	// Map gRPC codes to GraphQL errors
	switch st.Code() {
	case codes.OK:
		return nil

	case codes.InvalidArgument:
		return &gqlerror.Error{
			Message: st.Message(),
			Extensions: map[string]interface{}{
				"code": "BAD_REQUEST",
			},
		}

	case codes.NotFound:
		return &gqlerror.Error{
			Message: st.Message(),
			Extensions: map[string]interface{}{
				"code": "NOT_FOUND",
			},
		}

	case codes.AlreadyExists:
		return &gqlerror.Error{
			Message: st.Message(),
			Extensions: map[string]interface{}{
				"code": "ALREADY_EXISTS",
			},
		}

	case codes.PermissionDenied:
		return &gqlerror.Error{
			Message: "You don't have permission to perform this action",
			Extensions: map[string]interface{}{
				"code": "FORBIDDEN",
			},
		}

	case codes.Unauthenticated:
		return &gqlerror.Error{
			Message: "Authentication required",
			Extensions: map[string]interface{}{
				"code": "UNAUTHENTICATED",
			},
		}

	case codes.ResourceExhausted:
		return &gqlerror.Error{
			Message: "Rate limit exceeded. Please try again later",
			Extensions: map[string]interface{}{
				"code": "RATE_LIMIT_EXCEEDED",
			},
		}

	case codes.FailedPrecondition:
		return &gqlerror.Error{
			Message: st.Message(),
			Extensions: map[string]interface{}{
				"code": "PRECONDITION_FAILED",
			},
		}

	case codes.Aborted:
		return &gqlerror.Error{
			Message: "Operation was aborted. Please try again",
			Extensions: map[string]interface{}{
				"code": "ABORTED",
			},
		}

	case codes.OutOfRange:
		return &gqlerror.Error{
			Message: st.Message(),
			Extensions: map[string]interface{}{
				"code": "OUT_OF_RANGE",
			},
		}

	case codes.Unimplemented:
		return &gqlerror.Error{
			Message: "This feature is not yet implemented",
			Extensions: map[string]interface{}{
				"code": "NOT_IMPLEMENTED",
			},
		}

	case codes.Unavailable:
		return &gqlerror.Error{
			Message: "Service temporarily unavailable. Please try again later",
			Extensions: map[string]interface{}{
				"code": "SERVICE_UNAVAILABLE",
			},
		}

	case codes.DeadlineExceeded:
		return &gqlerror.Error{
			Message: "Request timeout. Please try again",
			Extensions: map[string]interface{}{
				"code": "TIMEOUT",
			},
		}

	case codes.Canceled:
		return &gqlerror.Error{
			Message: "Request was canceled",
			Extensions: map[string]interface{}{
				"code": "CANCELED",
			},
		}

	case codes.DataLoss:
		return &gqlerror.Error{
			Message: "Data loss or corruption detected",
			Extensions: map[string]interface{}{
				"code": "DATA_LOSS",
			},
		}

	case codes.Unknown, codes.Internal:
		fallthrough
	default:
		// Don't expose internal error details
		return &gqlerror.Error{
			Message: "An internal error occurred. Please try again later",
			Extensions: map[string]interface{}{
				"code": "INTERNAL_ERROR",
			},
		}
	}
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) error {
	return &gqlerror.Error{
		Message: message,
		Extensions: map[string]interface{}{
			"code": "BAD_REQUEST",
		},
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) error {
	return &gqlerror.Error{
		Message: resource + " not found",
		Extensions: map[string]interface{}{
			"code": "NOT_FOUND",
		},
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError() error {
	return &gqlerror.Error{
		Message: "Authentication required",
		Extensions: map[string]interface{}{
			"code": "UNAUTHENTICATED",
		},
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError() error {
	return &gqlerror.Error{
		Message: "You don't have permission to perform this action",
		Extensions: map[string]interface{}{
			"code": "FORBIDDEN",
		},
	}
}

// NewInternalError creates an internal error
func NewInternalError() error {
	return &gqlerror.Error{
		Message: "An internal error occurred. Please try again later",
		Extensions: map[string]interface{}{
			"code": "INTERNAL_ERROR",
		},
	}
}
