package handler

import (
	"fmt"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Backend is healthy!")
}

func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", h.Check)
}
