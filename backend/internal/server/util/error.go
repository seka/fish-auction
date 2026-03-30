package util

import (
	"errors"
	"log"
	"net/http"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
)

// ErrorResponse represents a JSON error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// HandleError converts domain errors to JSON responses.
func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var status int
	var errorType, message string

	var valErr *domainErrors.ValidationError
	var notFoundErr *domainErrors.NotFoundError
	var conflictErr *domainErrors.ConflictError
	var unauthErr *domainErrors.UnauthorizedError
	var forbiddenErr *domainErrors.ForbiddenError
	var goneErr *domainErrors.GoneError

	if errors.As(err, &valErr) {
		status = http.StatusBadRequest
		errorType = "validation_error"
		message = valErr.Error()
	} else if errors.As(err, &notFoundErr) {
		status = http.StatusNotFound
		errorType = "not_found"
		message = notFoundErr.Error()
	} else if errors.As(err, &conflictErr) {
		status = http.StatusConflict
		errorType = "conflict"
		message = conflictErr.Error()
	} else if errors.As(err, &unauthErr) {
		status = http.StatusUnauthorized
		errorType = "unauthorized"
		message = unauthErr.Error()
	} else if errors.As(err, &forbiddenErr) {
		status = http.StatusForbidden
		errorType = "forbidden"
		message = forbiddenErr.Error()
	} else if errors.As(err, &goneErr) {
		status = http.StatusGone
		errorType = "gone"
		message = goneErr.Error()
	} else {
		status = http.StatusInternalServerError
		errorType = "internal_error"
		message = "An internal error occurred"
		log.Printf("Internal error: %v", err)
	}

	resp := ErrorResponse{Error: errorType, Message: message, Code: status}
	WriteJSON(w, status, resp)
}

// WriteError writes an error response with the given status code and message.
func WriteError(w http.ResponseWriter, status int, message string) {
	resp := ErrorResponse{
		Error:   "error",
		Message: message,
		Code:    status,
	}
	WriteJSON(w, status, resp)
}
