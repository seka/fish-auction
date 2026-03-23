package errors

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
)

// HandleError converts database-specific errors (like sql.ErrNoRows or pq.Error)
// into domain-specific apperrors (like NotFoundError or ConflictError)
// resourceName is used for formatting NotFoundError
// resourceID is used for formatting NotFoundError
// errCtxMessage is used to prefix the original error when converted to a standard error.
func HandleError(err error, resourceName string, resourceID any, errCtxMessage string) error {
	if err == nil {
		return nil
	}

	// Handle standard SQL "Not Found" error
	if errors.Is(err, sql.ErrNoRows) {
		return &domainErrors.NotFoundError{Resource: resourceName, ID: resourceID}
	}

	// Handle PostgreSQL specific errors
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch string(pqErr.Code) {
		case pgerrcode.UniqueViolation:
			// Often occurs when an item with the same unique key already exists
			return &domainErrors.ConflictError{Message: fmt.Sprintf("%s already exists", resourceName)}
		case pgerrcode.ForeignKeyViolation:
			// Often occurs when referring to a non-existent parent
			// Extract constraint name or detail if possible for better messaging,
			// or default to a standard conflict/not-found message.
			// Here we map it to ConflictError to indicate a relational integrity issue.
			detail := pqErr.Detail
			if detail == "" {
				detail = pqErr.Message
			}
			return &domainErrors.ConflictError{Message: fmt.Sprintf("foreign key constraint failed on %s: %s", resourceName, detail)}
		}
	}

	// Check if the error is already a domain error, return as is
	var notFoundErr *domainErrors.NotFoundError
	var conflictErr *domainErrors.ConflictError
	var validationErr *domainErrors.ValidationError
	if errors.As(err, &notFoundErr) || errors.As(err, &conflictErr) || errors.As(err, &validationErr) {
		return err
	}

	// For unhandled internal errors, wrap with context
	if errCtxMessage != "" {
		// Provide context without creating a new line
		return fmt.Errorf("%s: %w", errCtxMessage, err)
	}

	// Fallback wrapping if no context was given
	return fmt.Errorf("database generic error for %s: %w", resourceName, err)
}

// IsUniqueViolation returns true if the error is a PostgreSQL unique violation.
// Can be used when more control over the custom ConflictError message is needed.
func IsUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && string(pqErr.Code) == pgerrcode.UniqueViolation
}
