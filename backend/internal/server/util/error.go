package util

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/server/dto"
)

// HandleError converts domain errors to JSON responses.
func HandleError(w http.ResponseWriter, err error) {
	var status int
	var errorType, message string
	switch e := err.(type) {
	case *errors.ValidationError:
		status = http.StatusBadRequest
		errorType = "validation_error"
		message = e.Error()
	case *errors.NotFoundError:
		status = http.StatusNotFound
		errorType = "not_found"
		message = e.Error()
	case *errors.ConflictError:
		status = http.StatusConflict
		errorType = "conflict"
		message = e.Error()
	case *errors.UnauthorizedError:
		status = http.StatusUnauthorized
		errorType = "unauthorized"
		message = e.Error()
	default:
		status = http.StatusInternalServerError
		errorType = "internal_error"
		message = "An internal error occurred"
		log.Printf("Internal error: %v", err)
	}
	resp := dto.ErrorResponse{Error: errorType, Message: message, Code: status}
	WriteJSON(w, status, resp)
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// WriteError writes an error response with the given status code and message.
func WriteError(w http.ResponseWriter, status int, message string) {
	resp := dto.ErrorResponse{
		Error:   "error",
		Message: message,
		Code:    status,
	}
	WriteJSON(w, status, resp)
}
