package handler

import (
	"fmt"
	"net/http"
)

// HealthHandler handles health check requests.
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check handles the health check request.
func (h *HealthHandler) Check(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Backend is healthy!")
}

// RegisterRoutes registers the health handler routes to the given mux.
func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", h.Check)
}
