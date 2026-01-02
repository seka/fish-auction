package postgres

import "github.com/lib/pq"

const ErrCodeUniqueViolation = pq.ErrorCode("23505")
const ErrCodeForeignKeyViolation = pq.ErrorCode("23503")
