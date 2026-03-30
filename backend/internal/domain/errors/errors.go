package errors

import "fmt"

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       any
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %v not found", e.Resource, e.ID)
}

// ConflictError represents a conflict error (e.g., item already sold)
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

// UnauthorizedError represents an authentication error
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

// ForbiddenError represents an authorization error
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

// GoneError represents a resource that is permanently gone
type GoneError struct {
	Resource string
	ID       any
	Message  string
}

func (e *GoneError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("%s with ID %v is gone", e.Resource, e.ID)
}
