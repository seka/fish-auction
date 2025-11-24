package handler

import (
	"encoding/json"
	"net/http"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Create(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if creds.Password == "admin-password" {
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_session",
			Value:    "authenticated",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   86400, // 1 day
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	} else {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
