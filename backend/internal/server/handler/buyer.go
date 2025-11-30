package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
)

type BuyerHandler struct {
	createUseCase buyer.CreateBuyerUseCase
	listUseCase   buyer.ListBuyersUseCase
	loginUseCase  buyer.LoginBuyerUseCase
}

func NewBuyerHandler(r registry.UseCase) *BuyerHandler {
	return &BuyerHandler{
		createUseCase: r.NewCreateBuyerUseCase(),
		listUseCase:   r.NewListBuyersUseCase(),
		loginUseCase:  r.NewLoginBuyerUseCase(),
	}
}

func (h *BuyerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBuyerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}
	buyer, err := h.createUseCase.Execute(r.Context(), req.Name, req.Email, req.Password, req.Organization, req.ContactInfo)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	resp := dto.BuyerResponse{ID: buyer.ID, Name: buyer.Name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BuyerHandler) List(w http.ResponseWriter, r *http.Request) {
	buyers, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		util.HandleError(w, err)
		return
	}
	resp := make([]dto.BuyerResponse, len(buyers))
	for i, b := range buyers {
		resp[i] = dto.BuyerResponse{ID: b.ID, Name: b.Name}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BuyerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginBuyerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	buyer, err := h.loginUseCase.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	// Set session cookie
	// In a real app, generate a secure token and store it in Redis/DB
	// For now, we'll use a simple signed token or just the ID (insecure but simple for MVP)
	// Let's use "buyer_session" cookie with value "authenticated" and maybe another for ID?
	// Or better, use a JWT or signed cookie.
	// Given the previous "admin_session" implementation, let's stick to simple cookie for now.
	// But we need the BuyerID for bidding.
	// So let's set "buyer_id" cookie as well (HttpOnly).

	http.SetCookie(w, &http.Cookie{
		Name:     "buyer_session",
		Value:    "authenticated",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
	})

	// Store Buyer ID in a separate cookie or encoded in the session
	// For simplicity, let's use a separate cookie "buyer_id"
	http.SetCookie(w, &http.Cookie{
		Name:     "buyer_id",
		Value:    fmt.Sprintf("%d", buyer.ID),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	resp := dto.BuyerResponse{ID: buyer.ID, Name: buyer.Name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BuyerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "buyer_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "buyer_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

func (h *BuyerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/buyers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			h.List(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/buyers/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/buyers/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Logout(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
