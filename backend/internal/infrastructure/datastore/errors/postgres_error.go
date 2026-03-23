package dserrors

import "github.com/lib/pq"

// ErrCodeUniqueViolation provides ErrCodeUniqueViolation related functionality.
const ErrCodeUniqueViolation = pq.ErrorCode("23505")
// ErrCodeForeignKeyViolation provides ErrCodeForeignKeyViolation related functionality.
const ErrCodeForeignKeyViolation = pq.ErrorCode("23503")
